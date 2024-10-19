package member

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andfxx27/gin-gonic-playground/config"
	"github.com/andfxx27/gin-gonic-playground/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func NewHandler(repo Repositorier, redisClient *redis.Client) Handler {
	return &handler{
		repo,
		redisClient,
	}
}

type Handler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)

	GetProfile(c *gin.Context)
}

type handler struct {
	repo        Repositorier
	redisClient *redis.Client
}

func (h *handler) SignUp(c *gin.Context) {
	errMsg := "member.Handler SignUp error, failed sign up."

	var request SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errs := util.ToErrResponse(err)
		log.Err(err).Msg(fmt.Sprintf("%s Failed to bind request body to struct.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"errors":     errs,
			"statusCode": ErrSignUpErrBindRequest,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to generate password hash.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to create member.",
			"statusCode": ErrSignUpErrHashPassword,
		})
		return
	}

	member := Member{
		ID:        uuid.New().String(),
		Username:  request.Username,
		Password:  string(hashedPassword),
		Email:     request.Email,
		CreatedAt: time.Now(),
	}
	createdMember, err := h.repo.CreateMember(&member, c)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to create member.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to create member.",
			"statusCode": ErrSignUpErrRepoCreateMember,
		})
		return
	}

	log.Info().Msg(fmt.Sprintf("Member ID: %s, created member successfully.", createdMember.ID))
	c.JSON(http.StatusOK, gin.H{
		"message":    "Member sign up OK",
		"statusCode": Success,
		"member":     createdMember,
	})
	return
}

func (h *handler) SignIn(c *gin.Context) {
	conf := config.NewConf()
	errMsg := "member.Handler SignIn error, failed sign in."

	var request SignInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errs := util.ToErrResponse(err)
		log.Err(err).Msg(fmt.Sprintf("%s Failed to bind request body to struct.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"errors":     errs,
			"statusCode": ErrSignInErrBindRequest,
		})
		return
	}

	member, err := h.repo.GetMemberByUsernameOrEmail(request.Identifier, c)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to get member by username or email.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to get member information, please check your credentials.",
			"statusCode": ErrSignInErrGetMemberInformation,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(request.Password))
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to compare password hash.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Please check your credentials.",
			"statusCode": ErrSignInErrCompareHashPassword,
		})
		return
	}

	baseHost := "http://localhost:%d"
	now := time.Now()
	accessToken, err := util.CreateJWTWithClaims(jwt.MapClaims{
		"exp": now.Add(time.Hour * 24).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"aud": fmt.Sprintf(baseHost, conf.Server.ClientPort),
		"iss": fmt.Sprintf(baseHost, conf.Server.Port),
		"sub": member.ID,
	})

	c.JSON(200, gin.H{
		"message":     "Member sign in OK",
		"statusCode":  Success,
		"accessToken": accessToken,
	})
}

func (h *handler) GetProfile(c *gin.Context) {
	errMsg := "member.Handler GetProfile error, failed to get member profile."

	subject, _ := c.Get("subject")
	log.Info().Msg(fmt.Sprintf("Getting member profile for member ID: %s.", subject))

	// Get member information from cache first
	redisKey := fmt.Sprintf("member-profile-%s", subject)
	val, err := h.redisClient.Get(c, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to get member profile.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to get member profile.",
			"statusCode": ErrGetProfileErrGetCachedMemberProfile,
		})
		return
	}
	if val != "" {
		cachedMember := Member{}
		err = json.Unmarshal([]byte(val), &cachedMember)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%s Failed to unmarshal member profile.", errMsg))
			c.JSON(http.StatusOK, gin.H{
				"error":      "Failed to get member profile.",
				"statusCode": ErrGetProfileErrUnmarshalMemberProfile,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Member get profile OK.",
			"statusCode": Success,
			"member":     cachedMember,
		})
		return
	}

	member, err := h.repo.GetMemberByID(fmt.Sprintf("%s", subject), c)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to get member profile.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to get member profile.",
			"statusCode": ErrGetProfileErrGetMemberProfile,
		})
		return
	}

	memberBytes, err := json.Marshal(member)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to marshal member profile into byte array.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to get member profile.",
			"statusCode": ErrGetProfileErrMarshalMemberProfile,
		})
	}

	err = h.redisClient.Set(c, redisKey, memberBytes, 0).Err()
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s Failed to set member profile to redis.", errMsg))
		c.JSON(http.StatusOK, gin.H{
			"error":      "Failed to get member profile.",
			"statusCode": ErrGetProfileErrSetMemberProfileToRedis,
		})
		return
	}

	c.JSON(200, gin.H{
		"message":    "Member get profile OK",
		"statusCode": Success,
		"member":     member,
	})
}
