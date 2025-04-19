package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/silasburger/lenslocked/context"
	"github.com/silasburger/lenslocked/errors"
	"github.com/silasburger/lenslocked/models"
)

type Users struct {
	Templates struct {
		New                Template
		SignIn             Template
		CurrentUser        Template
		ForgotPassword     Template
		CheckYourEmail     Template
		ResetPassword      Template
		PasswordlessSignin Template
		EditEmail          Template
	}
	UsersService         *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UsersService.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrPasswordMismatch) {
			err = errors.Public(err, "Incorrect email or password.")
		}
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}

	u.signInUser(w, r, user)
}

func (u Users) signInUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	var data struct {
		Email string
	}
	data.Email = user.Email
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		err = errors.Public(err, "Something went wrong.")
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UsersService.Create(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account.")
		}
		u.Templates.New.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning about not being able to sign
		// the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	u.Templates.CurrentUser.Execute(w, r, user)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) PasswordlessSignin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.PasswordlessSignin.Execute(w, r, data)
}

func (u Users) ProcessPasswordlessSignin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(models.ErrNotFound, "There is no account with that email address.")
		}
		u.Templates.PasswordlessSignin.Execute(w, r, data, err)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}

	fmt.Println("📙 📙 📙 📙 📙 📙 📙 📙 ", pwReset.Token, "reset token", "📙📙 📙 📙 📙 📙 📙 📙 📙 📙 📙 📙 ")

	//TODO: Make URL here configurable
	resetURL := "localhost:3000/email-signin?" + vals.Encode()

	fmt.Println("✨✨✨✨✨", resetURL, "✨✨✨✨✨✨✨✨✨✨✨")

	err = u.EmailService.PasswordlessSignin(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ProcessEmailSignin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	fmt.Println("😍😍😍😍😍😍😍😍😍😍😍", data.Token, "😍😍😍😍😍😍😍😍😍😍😍😍😍😍😍😍😍😍")
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	u.signInUser(w, r, user)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	// verify token with token service
	// replace password with new password
	// create a new session for them and sign them in
	// redirect them to the users/me page
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		// TODO: Distinguish between types of errors
		if errors.Is(err, models.ErrPasswordResetTokenInvalid) {
			http.Error(w, "Token invalid.", http.StatusBadRequest)
		}
		if errors.Is(err, models.ErrPasswordResetTokenExpired) {
			http.Error(w, "Token expired.", http.StatusBadRequest)
		}
		fmt.Println(err)
		return
	}
	err = u.UsersService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Sign the user in. At this point the password has already been reset
	// so if there is a problem signing them in we can simply put them at the sign-in page
	u.signInUser(w, r, user)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(models.ErrNotFound, "There is no account with that email address.")
		}
		u.Templates.ForgotPassword.Execute(w, r, data, err)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}

	//TODO: make URL configurable
	resetURL := "localhost:3000/reset-pw?" + vals.Encode()

	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) EditEmail(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	u.Templates.EditEmail.Execute(w, r, user)
}

func (u Users) ProcessEditEmail(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	user := context.User(r.Context())
	data.Email = r.FormValue("email")

	err := u.UsersService.UpdateEmail(user.ID, data.Email)
	if err != nil {
		err = errors.Public(err, "Something went wrong.")
		u.Templates.EditEmail.Execute(w, r, data, err)
		return
	}

	http.Redirect(w, r, "/users/me", http.StatusFound)
}
