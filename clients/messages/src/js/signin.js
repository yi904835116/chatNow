// add code here that is specific to the sign-in page (index.html)

"use strict";

var signUpForm = document.getElementById("signin-form");

signUpForm.addEventListener("submit", function (evt) {
    evt.preventDefault();
    handleSubmitSigninForm(evt);
});

// Sign in the user and save the session token in local storage.
function handleSubmitSigninForm(e) {
    e.preventDefault();

    var emailInput = document
        .getElementById("email-input")
        .value;
    var passwordInput = document
        .getElementById("password-input")
        .value;

    let credential = {
        email: emailInput,
        password: passwordInput
    };

    let url;

    url = 'https://api.patrick-yi.com/v1/sessions';

    fetch(url, {
        method: 'post',
        body: JSON.stringify(credential),
        mode: 'cors',
        headers: new Headers({'Content-Type': 'application/json'})
    }).then(res => {
        // If we get a successful response (status code < 300), save the contents of the
        // Authorization response header to local storage.
        if (res.status < 300) {
            // Save session token to local storage.
            const sessionToken = res
                .headers
                .get('Authorization');

            if (sessionToken != null) {
                localStorage.setItem('session-token', sessionToken);
            }
            return res.json();
        }

        // If response is not ok, catch the error contained in body.
        return res.text();
    }).then(data => {
        // If data type is string, it means this is an error sent by server.
        if (typeof data === 'string') {
            throw Error(data);
        } else {
            window
                .location
                .replace('app.html');
        }
    }).catch(error => {
        alert(error);
    });
}
