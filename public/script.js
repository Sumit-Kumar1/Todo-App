document.addEventListener("htmx:responseError", (evt) => {
    console.log(evt);
    document.getElementById("errors").innerHTML = evt.detail.xhr.response;
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