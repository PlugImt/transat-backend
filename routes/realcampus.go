package routes

import (
	"database/sql"

	"Transat_2.0_Backend/middlewares"
	"Transat_2.0_Backend/realcampus/posts"
	"github.com/gofiber/fiber/v2"
)

// SetupRealCampusRoutes configures all the routes for RealCampus feature
func SetupRealCampusRoutes(router fiber.Router, db *sql.DB) {
	// Create a group for RealCampus routes
	realcampus := router.Group("/realcampus")

	// Post-related routes
	realcampus.Get("/posts/today", middlewares.JWTMiddleware, posts.GetUserTodayPosts(db))
	realcampus.Post("/posts", middlewares.JWTMiddleware, posts.CreatePost(db))

	/*
			realcampus.Get("/posts/friends", handlers.GetFriendPosts(db))
			realcampus.Get("/posts/friends/today", handlers.GetFriendTodayPosts(db))

		// Post reactions
			realcampus.Get("/posts/:postID/reactions", handlers.GetPostReactions(db))
			realcampus.Post("/posts/:postID/reactions", handlers.AddReaction(db))

		// Friendship-related routes
			realcampus.Get("/friendships", handlers.GetUserFriendships(db))
			realcampus.Post("/friendships", handlers.AddFriend(db))
			realcampus.Put("/friendships/:friendshipID/respond", handlers.RespondToFriendRequest(db))
			realcampus.Delete("/friendships/:friendshipID", handlers.RemoveFriend(db))

	*/
}
