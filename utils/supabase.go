package utils

import (
	"log"
	"os"

	"github.com/nedpals/supabase-go"
)

var SupabaseClient *supabase.Client

// InitializeSupabase sets up the Supabase client
func InitializeSupabase() {
	// Load environment variables
	LoadEnv()

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseUrl == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL or SUPABASE_KEY is not set in environment variables")
	}

	// Create Supabase client
	SupabaseClient = supabase.CreateClient(supabaseUrl, supabaseKey)
	log.Println("Supabase client initialized successfully")
}
