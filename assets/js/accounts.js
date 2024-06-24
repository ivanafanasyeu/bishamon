const addAccountInitiator = document.getElementById("btn-add-account");
const addAccountDialog = document.getElementById("dialog-add-account");
const closeAddAccountDialogIcon = document.getElementById("close");

addAccountInitiator.addEventListener("click", () => {
	addAccountDialog.showModal();
});

closeAddAccountDialogIcon.addEventListener("click", () => {
	addAccountDialog.close();
});

// htmx send data as HTML, where checkbox is sending as "on" or "off", instead of boolean
document.body.addEventListener("htmx:configRequest", function (evt) {
	if (evt.target.getAttribute("data-form-type") == "form-account") {
		evt.detail.parameters["isInTotalBalance"] = evt.target.isInTotalBalance.checked ? true : false;
	}
});

document.body.addEventListener("accountCreated", function (evt) {
	if (evt.target.getAttribute("data-form-type") == "form-account") {
		evt.target.reset();
		evt.target.querySelector("#form-errors-add-account").innerHTML = "";
		addAccountDialog.close();
	}
});

document.body.addEventListener("accountUpdated", function (evt) {
	if (evt.target.getAttribute("data-form-type") == "form-account") {
		evt.target.reset();
		evt.target.querySelector("#form-errors-edit-account").innerHTML = "";
		document.getElementById("dialog-edit").close();
	}
});
