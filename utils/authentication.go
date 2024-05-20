package utils

import "net/http"

// CheckAdmin checks if the user is an admin and sends a forbidden response if not.
func CheckAdmin(w http.ResponseWriter, r *http.Request) bool {

	if !isAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}
	return true
}

// TODO do something with this function

// Check if the user is an admin
func isAdmin(r *http.Request) bool {
	// Implement your admin check logic here
	// This could involve checking a token, session, or user role in the request
	return true // Placeholder, replace with actual logic xD
}
