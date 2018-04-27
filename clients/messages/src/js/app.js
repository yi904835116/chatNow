var user = {};
var hasUser = false;

var firstName = document.querySelector("#firstName");
var lastName = document.querySelector("#lastName");

var button = document.querySelector("#signOut");

function checkRefresh() {
    authenticateUser();
}

function authenticateUser() {
    const sessionToken = this.getSessionToken();

    const url = `https://api.patrick-yi.com/v1/users/me`;
    // Validate this session token.
    fetch(url, {
        method: 'get',
        mode: 'cors',
        headers: new Headers({Authorization: sessionToken})
    }).then(res => {
        if (res.status < 300) {
            return res.json();
        }
        return res.text();
    }).then(data => {
        if (typeof data === 'string') {
            throw Error(data);
        } else {
            user = JSON.parse(data);
            hasUser = true;
            
            firstName.appendChild(user.FirstName)
            lastName.appendChild(user.LastName)
        }
    }).catch(error => {
        window.alert(error.message);
        window
            .location
            .replace('index.html');
    });
}

// Get session token from local storage.
function getSessionToken() {
    const sessionToken = localStorage.getItem('session-token');
    if (sessionToken == null || sessionToken.length === 0) {
        // If no session token found in local storage, redirect the user back to landing
        // page.
        window
            .location
            .replace('index.html');
    }
    return sessionToken;
}


button.addEventListener("click", signOut)

// Sign out the user and end the session.
function signOut(e) {
    e.preventDefault();

    const sessionToken = this.getSessionToken();
    const url = `https://api.patrick-yi.com/v1/sessions/mine`;

    fetch(url, {
        method: 'delete',
        mode: 'cors',
        headers: new Headers({Authorization: sessionToken})
    }).then(res => {
        // If the response is successful, remove session token in local storage.
        if (res.status < 300) {
            this.setState({hasUser: false});
            localStorage.removeItem('session-token');
            window
                .location
                .replace('index.html');
        }
        return res.text();
    }).then(data => {
        if (typeof data === 'string') {
            throw Error(data);
        }
    }).catch(error => {
        alert(error.message);
    });
}