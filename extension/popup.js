document
  .getElementById('form-product_create')
  .addEventListener('submit', async (e) => {

    e.preventDefault();
    const name = document.getElementById('product-name').value;
    const origin = document.getElementById('product-origin').value;

    // Get ingredient list
    const ingredients = Array.from(
      document.querySelectorAll(
        '#ingredient-tag-input'
      )
    ).map(i => i.value);

    if (!name) {
      alert('Name  required');
      return;
    }

    if (!ingredients) {
      alert('at least one ingredient required');
      return;
    }

    //const ingredient_list = ingredients.join(', ');

    try {

      // QUERY PARAM VERSION
      const params = new URLSearchParams({
        name,
        origin,
        ingredients
      });

      const response = await fetch(
        `https://cff4ha.godiegogo.me/product/create?${params.toString()}`,
        {
          method: 'PUT'
        }
      );

      if (!response.ok) {
        throw new Error('Failed request');
      }

      const data = await response.json();
      console.log('SUCCESS:', data);
      alert('Product created successfully');

    } catch (err) {
      console.error(err);
      alert('Error creating product');
    }
  });

document.getElementById('close-btn').addEventListener('click', () => {
  window.close();
});
