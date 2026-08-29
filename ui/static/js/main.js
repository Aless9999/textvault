document.addEventListener('DOMContentLoaded', function () {
	// Подсветка активной ссылки в nav
	var navLinks = document.querySelectorAll("nav a");
	for (var i = 0; i < navLinks.length; i++) {
		var link = navLinks[i];
		if (link.getAttribute('href') === window.location.pathname) {
			link.classList.add("live");
			break;
		}
	}

	// Логика модального окна
	const deleteForm = document.querySelector('.delete-form');
	const deleteDialog = document.getElementById('delete-dialog');
	const cancelBtn = document.getElementById('dialog-cancel-btn');
	const confirmBtn = document.getElementById('dialog-confirm-btn');

	if (deleteForm && deleteDialog && cancelBtn && confirmBtn) {
		// Открытие диалога при клике на Delete
		deleteForm.addEventListener('submit', function (event) {
			event.preventDefault();
			deleteDialog.showModal();
		});

		// Отмена
		cancelBtn.addEventListener('click', function (e) {
			e.preventDefault();
			deleteDialog.close();
		});

		// Подтверждение
		confirmBtn.addEventListener('click', function (e) {
			e.preventDefault();
			deleteDialog.close();
			deleteForm.submit();
		});
	}
});