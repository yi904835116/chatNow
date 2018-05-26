
"use strict";

var signUpForm = document.getElementById("signup-form");
signUpForm.addEventListener("submit", function(evt) {
    evt.preventDefault();
    handleSubmitSignupForm(evt);
});


function handleSubmitSignupForm(e) {
    var emailInput = document.getElementById("email").value;
    var passwordInput = document.getElementById("password-input").value;
    var passwordConfirmation = document.getElementById("passwordConfirmation").value;
    var userNameInput = document.getElementById("userName-input").value;
    var lastNameInput = document.getElementById("lastName-input").value;
    var firstNameInput = document.getElementById("firstName-input").value;

    let user = {
        userName: userNameInput,
        lastName: lastNameInput,
        firstName: firstNameInput,
        email: emailInput,
        password: passwordInput,
        passwordConf: passwordConfirmation
    };

    let url = 'https://api.patrick-yi.com/v1/users';

    fetch(url, {
        method: 'post',
        body: JSON.stringify(user),
        mode: 'cors',
        headers: new Headers({
            'Content-Type': 'application/json'
        })
    })
        .then(res => {
            if (res.status < 300) {
                // Save session token to local storage.
                const sessionToken = res.headers.get('Authorization');

                if (sessionToken != null) {
                    localStorage.setItem('session-token', sessionToken);
                }
                return res.json();
            }

            return res.text();
        })
        .then(data => {
            if (typeof data === 'string') {
                throw Error(data);
            } else {
                window.location.replace('app.html');
            }
        })
        .catch(error => {
            alert(error)
        });
}