package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"server/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func ResumeScan() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
		apiKey := os.Getenv("GEMINI_API_KEY")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var aiscan models.ResumeScan
		if err := c.BindJSON(&aiscan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validationErr := validate.Struct(aiscan)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		opts := genai.UploadFileOptions{DisplayName: "Gemini pdf"}
		f, err := os.Open(*aiscan.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		doc1, err := client.UploadFile(ctx, "", f, &opts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Uploaded file %s as: %q\n", doc1.DisplayName, doc1.URI)
		model := client.GenerativeModel("gemini-1.5-flash")

		prompt := []genai.Part{
			genai.FileData{URI: doc1.URI},
			genai.Text(*aiscan.Prompt + "for role" + *aiscan.Role),
		}

		resp, err := model.GenerateContent(ctx, prompt...)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, resp)
	}

}

func Summarize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.VideoIdRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		videoUrl := request.VideoUrl

		fmt.Printf("Received url %s\n", videoUrl)

		cmd := exec.Command("python", "script.py", videoUrl)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute Python script"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"output":    out.String(),
			"video_url": videoUrl,
		})

	}
}

func ResumeScore() gin.HandlerFunc {

	return func(c *gin.Context) {
		description := c.PostForm("description")
		file, _, err := c.Request.FormFile("resume")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
			return
		}

		file.Close()

		tmpFile, err := os.CreateTemp("", "resume_*.pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		}

		defer os.Remove(tmpFile.Name())
		io.Copy(tmpFile, file)

		cmd := exec.Command("python", "parse.py", tmpFile.Name(), description)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute Python script"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"match_percent": out.String(),
		})
	}
}
