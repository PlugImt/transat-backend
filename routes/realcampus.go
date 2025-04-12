package routes

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"Transat_2.0_Backend/realcampus/friendships"
	"Transat_2.0_Backend/realcampus/posts"
	"Transat_2.0_Backend/realcampus/users"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// SetupRealCampusRoutes configures all the routes for RealCampus feature
func SetupRealCampusRoutes(router fiber.Router, db *sql.DB) {
	// Create a group for RealCampus routes.
	realcampus := router.Group("/realcampus")

	// Post-related routes.
	// Group all posts routes under /api/posts.
	postsGroup := realcampus.Group("/posts")
	postsGroup.Get("/today", jwtMiddleware, posts.GetUserTodayPosts(db))
	postsGroup.Post("/", jwtMiddleware, posts.CreatePost(db))

	// Friendship-related routes.
	// Group all friendship routes under /api/friendships.
	friendshipsGroup := realcampus.Group("/friendships")
	friendshipsGroup.Get("/", jwtMiddleware, friendships.GetUserFriendships(db))
	friendshipsGroup.Post("/send", jwtMiddleware, friendships.SendFriendRequest(db))
	friendshipsGroup.Post("/reject", jwtMiddleware, friendships.RejectFriendRequest(db))
	friendshipsGroup.Post("/accept", jwtMiddleware, friendships.AcceptFriendRequest(db))
	friendshipsGroup.Post("/remove", jwtMiddleware, friendships.RemoveFriend(db))
	//friendshipsGroup.Post("/cancel", jwtMiddleware, friendships.CancelFriendRequest(db))

	router.Get("/users/search", jwtMiddleware, users.SearchUsers(db))

	/*
			realcampus.Get("/posts/friends", handlers.GetFriendPosts(db))
			realcampus.Get("/posts/friends/today", handlers.GetFriendTodayPosts(db))

		// Post reactions
			realcampus.Get("/posts/:postID/reactions", handlers.GetPostReactions(db))
			realcampus.Post("/posts/:postID/reactions", handlers.AddReaction(db))

			realcampus.Post("/friendships_group", handlers.AddFriend(db))
			realcampus.Put("/friendships_group/:friendshipID/respond", handlers.RespondToFriendRequest(db))
			realcampus.Delete("/friendships_group/:friendshipID", handlers.RemoveFriend(db))

	*/
}

/*
*
IMPORTANT : AS security.go SHOULD BE PLACED IN A DEDICATED PACKAGE AND THE ENTIRE PROJECT
ORGANIZATION SHOULD BE REFACTORED, THE CODE BELOW SHOULD BE REMOVED AFTER REFACTOR. WHILE
REFACTOR IS NOT DONE, CODE BELOW IS NECESSARY.
*/
func jwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	log.Println("╔======== 📧 JWT Middleware 📧 ========╗")

	if authHeader == "" {
		log.Println("║ 💥 Missing token")
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[7:]
	}

	token, err := validateJWT(tokenString)
	if err != nil {
		log.Println("║ 💥 Invalid token: ", err)
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("║ 💥 Invalid claims")
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	c.Locals("email", claims["email"])

	log.Println("║ ✅ Token is valid")
	log.Println("║ 📧 Email: ", claims["email"])
	log.Println("╚=======================================╝")

	return c.Next()
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
