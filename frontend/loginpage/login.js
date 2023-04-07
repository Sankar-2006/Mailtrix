const inputs = document.querySelectorAll('.input');
const mailContain = document.querySelector('.inptCont-one');
const passContain = document.querySelector('.inptCont-two');
const mailInput = document.querySelector('#emailInpt');
const passInput = document.querySelector('#passInpt');
const mailIcon = document.querySelector('.fa-envelope');
const keyIcon = document.querySelector('.fa-key');
const mailLabel = document.querySelector('.label-one');
const passLabel = document.querySelector('.label-two');
const errmsgMail = document.querySelector('#errorMail');
const errmsgPass = document.querySelector('#errorPass');
const button = document.querySelector('#button');

inputs.forEach(input => {
    input.addEventListener('focus', () => {
        const addClass = input.parentNode.parentNode;
        addClass.classList.add('focus');
    })
    input.addEventListener('blur', () => {
        const removeClass = input.parentNode.parentNode;
        if (input.value === "") {
            removeClass.classList.remove('focus');
        }
    })
})

/******************* ---VALIDATION--- ******************* */



function checkEmptyFun() {
    if (mailInput.value === "") {
        errmsgMail.innerHTML = "Email cannot be Empty";
        // mailContain.style.borderBottomColor = "red";
        // mailContain.style.transition = "none";
        mailIcon.style.color = "red";
        mailIcon.style.transition = "none";
        // mailLabel.style.color = "red"
    } else {
        checkCrtFormFun();
    }

    if (passInput.value === "") {
        errmsgPass.innerHTML = "Password cannot be Empty";
        // passContain.style.borderBottomColor = "red";
        // passContain.style.transition = "none";
        keyIcon.style.color = "red";
        keyIcon.style.transition = "none";
        // passLabel.style.color = "red"
    } else {
        checkCrtFormFun();
    }
}

function checkCrtFormFun() {

    let isCrtMail = '';
    let isCrtPass = true;

    if (isCrtMail != true) {
        errmsgMail.innerHTML = "Email is incorrect or does not Exists";
        // mailContain.style.borderBottomColor = "red";
        // mailContain.style.transition = "none";
        mailIcon.style.color = "red";
        mailIcon.style.transition = "none";
        // mailLabel.style.color = "red"
    }

    if (isCrtPass != true) {
        errmsgPass.innerHTML = "Password is Incorrect";
        // passContain.style.borderBottomColor = "red";
        // passContain.style.transition = "none";
        keyIcon.style.color = "red";
        keyIcon.style.transition = "none";
        // passLabel.style.color = "red"
    }
}

button.addEventListener('click', (e) => {
    e.preventDefault();
    checkEmptyFun();
})
