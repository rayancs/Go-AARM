package controller

import (
	"app/configs"
	"app/middleware"
	repo "app/repos"
	"app/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
)

type SampleTypeReq struct {
	Questions []SampleType
}
type SampleType struct {
	Question    string
	Option1     string
	Option2     string
	Option3     string
	Option4     string
	Explination string
}

func BasePage(g *gin.Engine) {
	g.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(HtmlDoc()))
	})
	g.GET("/api", middleware.JWTMiddleware(), func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(HtmlDoc()))
	})
	g.GET("/claims", middleware.JWTMiddleware(), func(ctx *gin.Context) {
		res, is := ctx.Get("claims")
		if !is {
			ctx.JSON(401, map[string]string{
				"missing": "jwt",
			})
			return
		}
		ctx.JSON(200, res)
	})
	g.GET("/rs", func(ctx *gin.Context) {
		ctx.JSON(200, configs.GenerateSchema[types.HttpResponseType]())
	})
	g.GET("/ai-test", func(ctx *gin.Context) {
		topic := ctx.Query("topic")
		if len(strings.Split(topic, " ")) > 1 {
			ctx.JSON(400, map[string]string{
				"error":     "only one topic allowed",
				"topic_len": fmt.Sprint(len(strings.Split(topic, " "))),
			})
			return
		}
		if len(topic) == 0 || topic == "" {
			ctx.JSON(400, map[string]string{
				"error": "empty topics not allowed",
			})
			return
		}
		key, _ := configs.GetOpenAiKeys()
		openAiRepo := repo.NewOpenAi(key.Key)

		msg := []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fmt.Sprintf("make 15 mcqs about ,  %s , also provide explination", topic)),
		}
		res, err := repo.StructuredText[SampleTypeReq](openAiRepo, msg, "mcq_schema")
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
		ctx.JSON(200, res)

	})

}
func SecureHTMLDoc() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Secure Page</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #0f172a;
            color: #e2e8f0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            overflow: hidden;
        }
        
        .container {
            text-align: center;
            position: relative;
            z-index: 10;
        }
        
        h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            position: relative;
            text-shadow: 0 0 10px #10b981;
        }
        
        .status {
            font-size: 1.5rem;
            margin-top: 2rem;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        .lock-container {
            width: 150px;
            height: 200px;
            margin: 20px auto;
            position: relative;
        }
        
        .lock-body {
            width: 100px;
            height: 80px;
            background-color: #10b981;
            border-radius: 10px;
            position: absolute;
            bottom: 0;
            left: 50%;
            transform: translateX(-50%);
            box-shadow: 0 0 20px rgba(16, 185, 129, 0.6);
        }
        
        .lock-shackle {
            width: 60px;
            height: 60px;
            border: 12px solid #10b981;
            border-bottom: none;
            border-radius: 50px 50px 0 0;
            position: absolute;
            top: 40px;
            left: 50%;
            transform: translateX(-50%);
            box-shadow: 0 0 20px rgba(16, 185, 129, 0.6);
            animation: secure 2s ease-in-out;
        }
        
        @keyframes secure {
            0% {
                transform: translateX(-50%) translateY(-50px);
            }
            50% {
                transform: translateX(-50%) translateY(-20px);
            }
            60% {
                transform: translateX(-50%) translateY(-25px);
            }
            70% {
                transform: translateX(-50%) translateY(-20px);
            }
            100% {
                transform: translateX(-50%) translateY(0);
            }
        }
        
        .checkmark {
            width: 30px;
            height: 30px;
            background-color: #10b981;
            border-radius: 50%;
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            display: flex;
            justify-content: center;
            align-items: center;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0% {
                transform: translate(-50%, -50%) scale(0.95);
                box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
            }
            
            70% {
                transform: translate(-50%, -50%) scale(1);
                box-shadow: 0 0 0 10px rgba(16, 185, 129, 0);
            }
            
            100% {
                transform: translate(-50%, -50%) scale(0.95);
                box-shadow: 0 0 0 0 rgba(16, 185, 129, 0);
            }
        }
        
        .check {
            width: 12px;
            height: 6px;
            border-left: 3px solid white;
            border-bottom: 3px solid white;
            transform: rotate(-45deg) translate(1px, -1px);
        }
        
        .shield {
            position: absolute;
            width: 40px;
            height: 40px;
            background-color: #10b981;
            border-radius: 50%;
            display: flex;
            justify-content: center;
            align-items: center;
            color: white;
            font-weight: bold;
            font-size: 20px;
            animation: shield-move 10s infinite linear;
            opacity: 0.8;
            z-index: -1;
        }
        
        @keyframes shield-move {
            0% {
                transform: translate(0, 0) rotate(0deg);
            }
            100% {
                transform: translate(100px, 100px) rotate(360deg);
            }
        }
        
        .security-grid {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            z-index: -2;
            perspective: 1000px;
            opacity: 0.3;
        }
        
        .grid-line {
            position: absolute;
            background-color: #10b981;
            box-shadow: 0 0 10px #10b981;
        }
        
        .horizontal {
            width: 100%;
            height: 1px;
            transform: rotateX(60deg);
        }
        
        .vertical {
            width: 1px;
            height: 100%;
            transform: rotateY(60deg);
        }
        
        .encryption-bits {
            position: absolute;
            width: 100%;
            height: 100%;
            z-index: -1;
        }
        
        .bit {
            position: absolute;
            color: #10b981;
            font-family: monospace;
            font-size: 12px;
            opacity: 0.8;
            animation: fall linear infinite;
            animation-duration: calc(3s + (var(--speed) * 5s));
            text-shadow: 0 0 5px #10b981;
        }
        
        @keyframes fall {
            0% {
                transform: translateY(-100vh);
            }
            100% {
                transform: translateY(100vh);
            }
        }
    </style>
</head>
<body>
    <div class="encryption-bits" id="encryption-bits"></div>
    <div class="security-grid" id="security-grid"></div>
    
    <div class="container">
        <h1>Secure Connection</h1>
        
        <div class="lock-container">
            <div class="lock-shackle"></div>
            <div class="lock-body">
                <div class="checkmark">
                    <div class="check"></div>
                </div>
            </div>
        </div>
        
        <div class="status">
            <span>This page is secure</span>
        </div>
    </div>
    
    <script>
        // Create security grid
        const securityGrid = document.getElementById('security-grid');
        const gridCount = 20;
        
        for (let i = 0; i < gridCount; i++) {
            const horizontalLine = document.createElement('div');
            horizontalLine.classList.add('grid-line', 'horizontal');
            horizontalLine.style.top = (100 / gridCount) * i + '%';
            securityGrid.appendChild(horizontalLine);
            
            const verticalLine = document.createElement('div');
            verticalLine.classList.add('grid-line', 'vertical');
            verticalLine.style.left = (100 / gridCount) * i + '%';
            securityGrid.appendChild(verticalLine);
        }
        
        // Create floating shield elements
        const container = document.querySelector('.container');
        for (let i = 0; i < 5; i++) {
            const shield = document.createElement('div');
            shield.classList.add('shield');
            shield.textContent = '🔒';
            shield.style.left = Math.random() * 100 + '%';
            shield.style.top = Math.random() * 100 + '%';
            shield.style.animationDuration = Math.random() * 20 + 10 + 's';
            shield.style.animationDelay = Math.random() * 5 + 's';
            document.body.appendChild(shield);
        }
        
        // Create binary/encryption bits animation
        const encryptionBits = document.getElementById('encryption-bits');
        const binaryChars = ['0', '1'];
        
        for (let i = 0; i < 100; i++) {
            const bit = document.createElement('div');
            bit.classList.add('bit');
            bit.style.left = Math.random() * 100 + '%';
            bit.style.setProperty('--speed', Math.random());
            bit.textContent = binaryChars[Math.floor(Math.random() * binaryChars.length)];
            encryptionBits.appendChild(bit);
        }
        
        // Update bits to create matrix-like effect
        setInterval(() => {
            const bits = document.querySelectorAll('.bit');
            bits.forEach(bit => {
                if (Math.random() > 0.9) {
                    bit.textContent = binaryChars[Math.floor(Math.random() * binaryChars.length)];
                }
            });
        }, 200);
    </script>
</body>
</html>`
}
func HtmlDoc() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Server Status</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #0f172a;
            color: #e2e8f0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            overflow: hidden;
        }
        
        .container {
            text-align: center;
            position: relative;
        }
        
        h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            position: relative;
            text-shadow: 0 0 10px #38bdf8;
        }
        
        .status {
            font-size: 1.5rem;
            margin-top: 2rem;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        .pulse {
            width: 20px;
            height: 20px;
            background-color: #4ade80;
            border-radius: 50%;
            margin-right: 15px;
            position: relative;
            animation: pulse 1.5s infinite;
        }
        
        @keyframes pulse {
            0% {
                transform: scale(0.95);
                box-shadow: 0 0 0 0 rgba(74, 222, 128, 0.7);
            }
            
            70% {
                transform: scale(1);
                box-shadow: 0 0 0 15px rgba(74, 222, 128, 0);
            }
            
            100% {
                transform: scale(0.95);
                box-shadow: 0 0 0 0 rgba(74, 222, 128, 0);
            }
        }
        
        .server {
            width: 100px;
            height: 120px;
            background-color: #475569;
            border-radius: 10px;
            margin: 0 auto;
            position: relative;
            overflow: hidden;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
        }
        
        .server-lights {
            display: flex;
            justify-content: space-around;
            padding: 10px;
        }
        
        .light {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #ef4444;
            animation: blink 2s infinite;
        }
        
        .light:nth-child(2) {
            animation-delay: 0.5s;
            background-color: #eab308;
        }
        
        .light:nth-child(3) {
            animation-delay: 1s;
            background-color: #4ade80;
        }
        
        @keyframes blink {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.3; }
        }
        
        .data-stream {
            position: absolute;
            width: 100%;
            height: 70px;
            bottom: 0;
            background-color: #1e293b;
            overflow: hidden;
        }
        
        .data {
            position: absolute;
            width: 2px;
            height: 10px;
            background-color: #38bdf8;
            opacity: 0.8;
            animation: data-flow 2s linear infinite;
        }
        
        @keyframes data-flow {
            0% {
                transform: translateY(-10px);
            }
            100% {
                transform: translateY(70px);
            }
        }
        
        .particles {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            z-index: -1;
        }
        
        .particle {
            position: absolute;
            background-color: #38bdf8;
            border-radius: 50%;
            opacity: 0.5;
            animation: float 15s infinite linear;
        }
        
        @keyframes float {
            0% {
                transform: translateY(0) translateX(0) rotate(0deg);
            }
            100% {
                transform: translateY(-100vh) translateX(100px) rotate(360deg);
            }
        }
    </style>
</head>
<body>
    <div class="particles" id="particles"></div>
    
    <div class="container">
        <h1>Server Status</h1>
        
        <div class="server">
            <div class="server-lights">
                <div class="light"></div>
                <div class="light"></div>
                <div class="light"></div>
            </div>
            <div class="data-stream" id="data-stream"></div>
        </div>
        
        <div class="status">
            <div class="pulse"></div>
            <span>Server is Running</span>
        </div>
    </div>
    
    <script>
        // Create data stream animation
        const dataStream = document.getElementById('data-stream');
        for (let i = 0; i < 30; i++) {
            const data = document.createElement('div');
            data.classList.add('data');
            data.style.left = Math.random() * 100 + '%';
            data.style.height = Math.random() * 15 + 5 + 'px';
            data.style.animationDuration = Math.random() * 3 + 1 + 's';
            data.style.animationDelay = Math.random() * 2 + 's';
            dataStream.appendChild(data);
        }
        
        // Create floating particles
        const particles = document.getElementById('particles');
        for (let i = 0; i < 20; i++) {
            const particle = document.createElement('div');
            particle.classList.add('particle');
            const size = Math.random() * 10 + 3;
            particle.style.width = size + 'px';
            particle.style.height = size + 'px';
            particle.style.left = Math.random() * 100 + '%';
            particle.style.top = Math.random() * 100 + '%';
            particle.style.animationDuration = Math.random() * 20 + 10 + 's';
            particle.style.animationDelay = Math.random() * 5 + 's';
            particles.appendChild(particle);
        }
    </script>
</body>
</html>
`
}
