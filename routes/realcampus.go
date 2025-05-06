package routes

import (
	"database/sql"

	"Transat_2.0_Backend/middlewares"
	"Transat_2.0_Backend/realcampus/friendships"
	"Transat_2.0_Backend/realcampus/posts"
	"Transat_2.0_Backend/realcampus/users"
	"Transat_2.0_Backend/utils"
	"github.com/gofiber/fiber/v2"
)

// SetupRealCampusRoutes configures all the routes for RealCampus feature
func SetupRealCampusRoutes(router fiber.Router, db *sql.DB) {
	// Create a group for RealCampus routes.
	realcampus := router.Group("/realcampus", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	// Post-related routes.
	// Group all posts routes under /api/posts.
	postsGroup := realcampus.Group("/posts", middlewares.JWTMiddleware)
	postsGroup.Get("/today", posts.GetUserTodayPosts(db))
	postsGroup.Post("/", posts.CreatePost(db))

	// Friendship-related routes.
	// Group all friendship routes under /api/friendships.
	friendshipsGroup := realcampus.Group("/friendships", middlewares.JWTMiddleware)
	friendshipsGroup.Get("/", friendships.GetUserFriendships(db))
	friendshipsGroup.Post("/send", friendships.SendFriendRequest(db))
	friendshipsGroup.Post("/reject", friendships.RejectFriendRequest(db))
	friendshipsGroup.Post("/accept", friendships.AcceptFriendRequest(db))
	friendshipsGroup.Post("/remove", friendships.RemoveFriend(db))
	friendshipsGroup.Post("/cancel", friendships.CancelFriendRequest(db))

	router.Get("/users/search", middlewares.JWTMiddleware, users.SearchUsers(db))

	/*
			realcampus.Get("/posts/friends", handlers.GetFriendPosts(db))
			realcampus.Get("/posts/friends/today", handlers.GetFriendTodayPosts(db))

		// Post reactions
			realcampus.Get("/posts/:postID/reactions", handlers.GetPostReactions(db))
			realcampus.Post("/posts/:postID/reactions", handlers.AddReaction(db))
	*/
}
