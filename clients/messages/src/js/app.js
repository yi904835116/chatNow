
var hasUser = false;

const signOutButton = document.querySelector("#signOut");
const firstName =document.querySelector("#firstName");
const lastName = document.querySelector("#lastName");
const searchInput =  document.querySelector("#searchUserInput");
const searchButton =  document.querySelector("#searchUserButton");
const searchList =  document.querySelector("#showSearchUsers");


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
            document.querySelector("#demo")
            

            hasUser = true;
            
            var firstNameText = document.createTextNode(data.firstName);
            var lastNameText = document.createTextNode(data.lastName);



            firstName.appendChild(firstNameText);
            lastName.appendChild(lastNameText);
        }
    }).catch(error => {
        window.alert(error.message);
        window
            .location
            .replace('index.html');
    });
}



searchButton.addEventListener("click", (e) =>{
    e.preventDefault();
    searchUser(e);
})


signOutButton.addEventListener("click", (e) =>{
    e.preventDefault();
    signOut(e)
})


function searchUser(e){
    const sessionToken = this.getSessionToken();

    const query = searchInput.value; //.trim().toLowerCase();
    const url = 'https://api.patrick-yi.com/v1/users?q=' + query;


    fetch(url, {
        method: 'get',
        mode: 'cors',
        headers: new Headers({
            Authorization: sessionToken
        })
    })
        .then(res => {
            if (res.status < 300) {
                return res.json();
            }
            return res.text();
        })
        .then(data => {
            if (typeof data === 'string') {
                throw Error(data);
            }        
            
            data.forEach(element => {

                let eachDiv = document.createElement("Div");
                let tempUserName = document.createTextNode("user name:  " + element.userName + "  ")
                let tempFirstName = document.createTextNode("first name:  " +element.firstName + "  ")
                let tempLastName = document.createTextNode("last name:  " +element.lastName + "  ")

                let tempImg = document.createElement("Img")
                tempImg.setAttribute("src","https://www.gravatar.com/avatar/6f34aaaed9dda84ac63618365fd6cd3c")

                eachDiv.appendChild(tempUserName);
                eachDiv.appendChild(tempFirstName);
                eachDiv.appendChild(tempLastName);
                eachDiv.appendChild(tempImg);

                searchList.appendChild(eachDiv);
            });


        })
        .catch(error => {
            alert(error);
        });

}



// Sign out the user and end the session.
function signOut(e) {

    const sessionToken = this.getSessionToken();
    const url = `https://api.patrick-yi.com/v1/sessions/mine`;

    fetch(url, {
        method: 'delete',
        mode: 'cors',
        headers: new Headers({Authorization: sessionToken})
    }).then(res => {
        // If the response is successful, remove session token in local storage.
        if (res.status < 300) {
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
