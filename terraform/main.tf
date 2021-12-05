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
    url     = "https://github.com/duyquang6/quote-today/archive/refs/tags/v0.0.2.tar.gz"
    version = "0.0.2"
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
    url = "https://${heroku_app.quote-today-app.name}.herokuapp.com/webhook/telebot"
  }
}

# Add a web-hook for the app
resource "heroku_app_webhook" "quote-today-app-release" {
  app_id  = heroku_app.quote-today-app.id
  level   = "notify"
  url     = "https://${heroku_app.quote-today-app.name}.herokuapp.com/webhook/telebot"
  include = ["api:release"]
}

output "quote_today_url" {
  value = "https://${heroku_app.quote-today-app.name}.herokuapp.com"
}
