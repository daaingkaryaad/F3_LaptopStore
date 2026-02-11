document.querySelectorAll('.category-card').forEach(card => {
  card.addEventListener('click', () => {
    alert('You selected: ' + card.querySelector('h2').innerText);
  });
});

document.querySelectorAll('.price-filters button').forEach(button => {
  button.addEventListener('click', () => {
    alert('Filtering laptops under: ' + button.innerText);
  });
});
