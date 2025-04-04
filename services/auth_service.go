package services

import (
	"app/configs"
	"app/logger"
	repo "app/repos"
	"app/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	AUTH_FLOW          = "authorization_code"
	TOKEN_EXCHANGE_API = "https://oauth2.googleapis.com/token"
	ACCESS_API         = "https://www.googleapis.com/oauth2/v2/userinfo"
	STANDARD_EMOJI     = "👻"
)

// login / signup same
type IAuthService interface {
	AuthUser() types.HttpResponseType
	GoogleSSO(code string) types.HttpResponseType
}

type AuthService struct {
	repo repo.IUserRepo
}

func NewAuthService(u repo.IUserRepo) *AuthService {
	return &AuthService{
		repo: u,
	}
}

func (a *AuthService) AuthUser() types.HttpResponseType {
	log := logger.New()
	err, _ := a.repo.LoginUser()
	if err != nil {
		log.Error(err.Error())
		return types.NewHttpResponse(
			configs.PleaseTryLater,
			"error-repo-auth",
			400,
			nil,
		)
	}
	return types.NewHttpResponse(configs.WelcomeBack, configs.Null, http.StatusOK, nil)
}

// GOOGLE AUTH SERVICE DOCS
//
//	For testing visit
//	"https://accounts.google.com/o/oauth2/auth?client_id=<your_client_id>&redirect_uri=<your_redirect_url>&response_type=code&scope=email profile"
//	1-after you authenticate with google it will redirect back to your website ( the redirect url , but with the query params ?code=<auth code>)
//	2-you can use the auth code to get a token in the server
//	3-the token can be used to get user data from gcps server
//	!!!The redirect url
//	http://localhost:3000/hub?code=4%2F0AQSTgQH7Y5VUULZVGM7m89jtJUQ3sZHFaZFnvAcNaJQz1VmYXHufoIxyqWGOgUWzcwZqlA&scope=email+profile+openid+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&authuser=0&prompt=consent
func (a *AuthService) GoogleSSO(code string) types.HttpResponseType {
	// logger instance
	log := logger.New()
	//	get server token
	data, err := a.exchangeForToken(code)
	if err != nil {
		log.Error(err.Error(), err)
		return types.NewHttpResponse("Could Not Log You In", "google_auth_token_exchange", 400, nil)
	}

	IAccessToken := data["access_token"]
	if IAccessToken == nil {
		log.Error("Access Token From Map Null")
		return types.NewHttpResponse("Could Not Log You In", "google_auth_token_exchange_null", 400, nil)
	}
	accessToken := fmt.Sprint(IAccessToken)

	userInfoMap, err := a.getUserInfo(accessToken)
	if err != nil {
		log.Error(err.Error(), err)
		return types.NewHttpResponse(configs.PleaseTryLater, "error_no_map_info", 400, nil)
	}
	userInfo, err := a.ExtractInfoFromMap(userInfoMap)
	if err != nil {
		log.Error(err.Error(), err)
		return types.NewHttpResponse(configs.PleaseTryLater, "error_no_map_info", 400, nil)
	}
	// search to find if user exsits or not
	res, err := a.repo.FindUser(userInfo.email)
	if err != nil {
		log.Error(err.Error())
		return types.NewHttpResponse(configs.PleaseTryLater, "error_at_search", 400, nil)
	}
	// check if no user exsists , if so create a new user
	if res == nil {
		// insert doc
		insertUserDoc := types.UserRequestType{
			Email:      userInfo.email,
			Emoji:      STANDARD_EMOJI,
			Name:       userInfo.givenName,
			PictureURL: userInfo.pictureUrl,
		}
		res, err := a.repo.CreateUser(insertUserDoc)
		if err != nil {
			log.Error(err.Error())
			return types.NewHttpResponse(configs.Oops, "bad_insert_new_user", http.StatusInternalServerError, nil)
		}
		// returns a id
		// create a jwt token
		token, err := CreateToken(res, insertUserDoc.Name, insertUserDoc.Email, insertUserDoc.Email)
		if err != nil {
			log.Error(err.Error())
			return types.NewHttpResponse(configs.Oops, "bad_token_gen_new_user", http.StatusInternalServerError, nil)
		}
		return types.NewHttpResponse(configs.WelcomeBack, configs.Null, 200, map[string]string{"token": token})
	}
	// for older users 
	token, err := CreateToken(res.Id, res.Name, res.Email, res.Email)
	if err != nil {
		log.Error(err.Error())
		return types.NewHttpResponse(configs.Oops, "bad_token_gen_new_user", http.StatusInternalServerError, nil)
	}
	return types.NewHttpResponse(configs.WelcomeBack, configs.Null, 200, map[string]string{"token": token})
}

// using the code from the client side you can exchange a token
//
//	you need the following paramaters :
//
//	{
//	 "code": "AUTH_CODE",
//	 "client_id": "YOUR_CLIENT_ID",
//	 "client_secret": "YOUR_CLIENT_SECRET",
//	 "redirect_uri": "YOUR_REDIRECT_URI",
//	 "grant_type": "authorization_code"
//	}
func (a *AuthService) exchangeForToken(code string) (map[string]interface{}, error) {
	//decode the code , as it is double encoded
	decodedCode, err := url.QueryUnescape(code)
	if err != nil {
		return map[string]interface{}{}, err
	}
	// read from envs
	gcpCreds, err := configs.GetGCPCreds()
	if err != nil {
		return map[string]interface{}{}, err
	}
	// creating the payload
	paramaters := url.Values{}
	paramaters.Add("code", decodedCode)
	paramaters.Add("client_id", gcpCreds.ClientID)
	paramaters.Add("client_secret", gcpCreds.ClientSecret)
	paramaters.Add("redirect_uri", gcpCreds.RedirectURI)
	paramaters.Add("grant_type", AUTH_FLOW)
	// post to gcp api
	res, err := http.PostForm(TOKEN_EXCHANGE_API, paramaters)
	// check errs
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer res.Body.Close()
	var tokenData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&tokenData); err != nil {
		return map[string]interface{}{}, err
	}
	return tokenData, nil

}

func (a *AuthService) getUserInfo(accessToken string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", ACCESS_API, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&userInfo)
	return userInfo, nil
}

type authInfo struct {
	email      string //`json:"email"`
	givenName  string //`name`
	pictureUrl string //`pictureUrl`
}

func (a *AuthService) ExtractInfoFromMap(data map[string]interface{}) (*authInfo, error) {
	// extract all info
	emailUntyped := data["email"]
	nameUntyped := data["given_name"]
	pictureUntyped := data["picture"]
	if emailUntyped == nil {
		return nil, errors.New("email_null")
	} else if nameUntyped == nil {
		return nil, errors.New("name_null")
	} else if pictureUntyped == nil {
		return nil, errors.New("name_null")
	}
	// convert to string
	email := fmt.Sprint(emailUntyped)
	name := fmt.Sprint(nameUntyped)
	picture := fmt.Sprint(pictureUntyped)
	//
	return &authInfo{
		email:      email,
		givenName:  name,
		pictureUrl: picture,
	}, nil

}
