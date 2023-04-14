const btnparent = document.querySelectorAll('.asideBtn');

btnparent.forEach((parent) => {
    parent.addEventListener(('click'), () => {
        btnparent.forEach((parent) =>
            parent.classList.remove('asideBtnActive')
        )
        parent.classList.add('asideBtnActive');
    })
    parent.addEventListener(('click'), () => {
        const sidebar = parent.dataset.forTab;
        const tabActive = document.querySelector(`.section[data-tab = "${sidebar}"]`);

        document.querySelectorAll('.section').forEach(tab => {
            tab.classList.remove('activesec')
            tabActive.classList.add('activesec')
        })
        
    })
})

/*====================================MailBox================================*/

const mailSection = document.querySelector('.inboxSection');
/**
 * note, this is just for testing purpose, soon the copmpose btn will be active *
 */
const newTest = document.querySelector('.composeBtn');

let lengthOfMails = document.querySelectorAll('.mailBox').length;

newTest.addEventListener(('click'), () => {
    lengthOfMails ++
    mailSection.innerHTML += `<section class='mailBox'>HI${lengthOfMails}</section>`;
})

