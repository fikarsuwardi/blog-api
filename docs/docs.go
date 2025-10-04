package docs

const SwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Blog API",
    "description": "REST API untuk Blog dengan Auth, Posts, dan Comments",
    "version": "1.0.0",
    "contact": {
      "name": "API Support"
    }
  },
  "host": "localhost:8080",
  "basePath": "/api",
  "schemes": ["http"],
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "Format: Bearer {token}"
    }
  },
  "paths": {
    "/register": {
      "post": {
        "tags": ["Auth"],
        "summary": "Register user baru",
        "parameters": [{
          "in": "body",
          "name": "body",
          "required": true,
          "schema": {
            "type": "object",
            "properties": {
              "email": {"type": "string", "example": "user@example.com"},
              "password": {"type": "string", "example": "password123"},
              "name": {"type": "string", "example": "John Doe"}
            }
          }
        }],
        "responses": {
          "201": {"description": "User created successfully"}
        }
      }
    },
    "/login": {
      "post": {
        "tags": ["Auth"],
        "summary": "Login user",
        "parameters": [{
          "in": "body",
          "name": "body",
          "required": true,
          "schema": {
            "type": "object",
            "properties": {
              "email": {"type": "string", "example": "user@example.com"},
              "password": {"type": "string", "example": "password123"}
            }
          }
        }],
        "responses": {
          "200": {"description": "Login successful"}
        }
      }
    },
    "/posts": {
      "get": {
        "tags": ["Posts"],
        "summary": "Get all posts",
        "responses": {
          "200": {"description": "List of posts"}
        }
      },
      "post": {
        "tags": ["Posts"],
        "summary": "Create new post",
        "security": [{"BearerAuth": []}],
        "parameters": [{
          "in": "body",
          "name": "body",
          "required": true,
          "schema": {
            "type": "object",
            "properties": {
              "title": {"type": "string", "example": "My First Post"},
              "content": {"type": "string", "example": "This is the content"}
            }
          }
        }],
        "responses": {
          "201": {"description": "Post created"}
        }
      }
    },
    "/posts/{id}": {
      "get": {
        "tags": ["Posts"],
        "summary": "Get post by ID",
        "parameters": [{
          "in": "path",
          "name": "id",
          "required": true,
          "type": "integer"
        }],
        "responses": {
          "200": {"description": "Post details"}
        }
      },
      "put": {
        "tags": ["Posts"],
        "summary": "Update post",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "title": {"type": "string"},
                "content": {"type": "string"}
              }
            }
          }
        ],
        "responses": {
          "200": {"description": "Post updated"}
        }
      },
      "delete": {
        "tags": ["Posts"],
        "summary": "Delete post",
        "security": [{"BearerAuth": []}],
        "parameters": [{
          "in": "path",
          "name": "id",
          "required": true,
          "type": "integer"
        }],
        "responses": {
          "200": {"description": "Post deleted"}
        }
      }
    },
    "/posts/{post_id}/comments": {
      "get": {
        "tags": ["Comments"],
        "summary": "Get comments for post",
        "parameters": [{
          "in": "path",
          "name": "post_id",
          "required": true,
          "type": "integer"
        }],
        "responses": {
          "200": {"description": "List of comments"}
        }
      },
      "post": {
        "tags": ["Comments"],
        "summary": "Create comment",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "in": "path",
            "name": "post_id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "content": {"type": "string", "example": "Great post!"}
              }
            }
          }
        ],
        "responses": {
          "201": {"description": "Comment created"}
        }
      }
    },
    "/posts/{post_id}/comments/{comment_id}": {
      "delete": {
        "tags": ["Comments"],
        "summary": "Delete comment",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "in": "path",
            "name": "post_id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "path",
            "name": "comment_id",
            "required": true,
            "type": "integer"
          }
        ],
        "responses": {
          "200": {"description": "Comment deleted"}
        }
      }
    }
  }
}`
