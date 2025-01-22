// Add a click event on buttons to open a specific modal
document.querySelectorAll('.modal-trigger')?.forEach(($trigger) => {
    const modal = $trigger.dataset.target;
    const $target = document.getElementById(modal);

    $trigger.addEventListener('click', () => {
        $target.classList.add('is-active')
    });
});

// Add a click event on various child elements to close the parent modal
document.querySelectorAll('.modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot .button')?.forEach(($close) => {
    const $target = $close.closest('.modal');

    $close.addEventListener('click', () => {
        $target.classList.remove('is-active')
    });
});

// Add a keyboard event to close all modals
document.addEventListener('keydown', (event) => {
    if (event.code == 'Escape') {
        document.querySelectorAll('.modal')?.forEach(($modal) => {
            $modal.classList.remove('is-active')
        })
    }
});
