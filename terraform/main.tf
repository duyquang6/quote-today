terraform {
  required_providers {
    heroku = {
      source  = "heroku/heroku"
      version = "~> 4.0"
    }
  }
}

variable "quote_today_app_name" {
  description = "Your app name"
}

resource "heroku_app" "quote-today-app" {
  name   = var.quote_today_app_name
  region = "us"
}

# Build code & release to the app
resource "heroku_build" "quote-today-app" {
  app        = heroku_app.quote-today-app.name
  buildpacks = ["https://github.com/heroku/heroku-buildpack-go.git"]

  source {
    path     = "/Users/linguyen/go/src/new-repo/part3"
    version = "0.0.1"
  }
}

# Launch the app's web process by scaling-up
resource "heroku_formation" "quote-today-app" {
  app        = heroku_app.quote-today-app.name
  type       = "web"
  quantity   = 1
  size       = "free"
  depends_on = [heroku_build.quote-today-app]
}

resource "heroku_addon" "pg-database" {
  app  = heroku_app.quote-today-app.name
  plan = "heroku-postgresql:hobby-dev"
}

# Add a web-hook addon for the app
resource "heroku_addon" "webhook" {
  app  = heroku_app.quote-today-app.name
  plan = "deployhooks:http"

  config = {
    url = "https://webhook.site/3d82d978-e781-4273-b7af-2cd717bacc44"
  }
}

output "quote_today_url" {
  value = "https://${heroku_app.quote-today-app.name}.herokuapp.com"
}
