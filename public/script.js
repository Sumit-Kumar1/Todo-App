document.addEventListener("htmx:responseError", (evt) => {
    const divErr = document.getElementById("errors")

    divErr.innerHTML = evt.detail.xhr.response; // show the error we got from backend

    setTimeout(function () { // remove everything inside divErr after 3 sec
        divErr.removeChild(divErr.childNodes[0]);
    }, 3000)
});


function updateModal(id) {
    let updateId = "update_" + id
    let updates = document.getElementById(updateId)
    updates.showModal()
}

function togglePasswordVisibility() {
    const eyeIcon = document.getElementById('toggleEye');
    const passwordInput = document.getElementById('password');
    if (passwordInput.type === 'password') {
        passwordInput.type = 'text';
        eyeIcon.classList.remove("fa-eye");
        eyeIcon.classList.add("fa-eye-slash");
    } else {
        passwordInput.type = 'password';
        eyeIcon.classList.remove("fa-eye-slash");
        eyeIcon.classList.add("fa-eye");
    }
}

function removeToast() {
    const toast = document.getElementById("notifi");
    if (toast == null) {
        return
    }

    toast.classList.add("hidden");
}