document.querySelectorAll('.file.has-name').forEach($file => {
    const $input = $file.querySelector('.file-input')
    $input.addEventListener('change', () => {
        $file.querySelector('.file-name').textContent = $input.files[0].name;
    })
})
