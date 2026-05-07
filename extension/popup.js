// function handleIngredientInput(e) {
//     const raw = e.target.value;
//     // Only split on ", " (comma + spce) to preserve names like "1,2-hexanediol"
//     const parts = raw.split(', ');
//     if (parts.length > 1) {
//       parts.slice(0, -1).forEach(p => commitIngredient(p.trim()));
//       commitIngredient(parts[parts.length - 1].trim());
//       e.target.value = "";
//     }
//   }

//   function handleIngredientKeydown(e) {
//     const input = e.target;
//     if ((e.key === 'Enter' || e.key === 'Tab') && input.value.trim()) {
//       e.preventDefault();
//       input.value.split(', ').forEach(p => commitIngredient(p.trim()));
//       input.value = '';
//     }
//     // Backspace on empty input removes the last tag
//     if (e.key === 'Backspace' && input.value === '') {
//       const box = document.getElementById('ingredient-tag-box');
//       const lastTag = box.querySelector('.ingredient-tag:last-of-type');
//       if (lastTag) removeIngredient(lastTag.dataset.name, lastTag);
//     }
//   }

//   function commitIngredient(name) {
//     if (!name) return;
//     const box = document.getElementById('ingredient-tag-box');
//     const hiddenContainer = document.getElementById('ingredient-hidden-inputs');
//     const input = document.getElementById('ingredient-tag-input');

//     // Deduplicate
//     if (box.querySelector(`.ingredient-tag[data-name="${CSS.escape(name)}"]`)) return;

//     const tag = document.createElement('span');
//     tag.className = 'ingredient-tag';
//     tag.dataset.name = name;
//     // tag.innerHTML = `${name}<button type="button" class="ingredient-tag-remove" aria-label="Remove ${name}" onclick="removeIngredient('${name.replace(/'/g, "\\'")}', this.parentElement)"><i class="bi bi-x"></i></button>`;
//     const text = document.createTextNode(name);

//     const btn = document.createElement('button');
//     btn.type = 'button';
//     btn.className = 'ingredient-tag-remove';
//     btn.innerHTML = '&times;';

//     btn.addEventListener('click', () => {
//       removeIngredient(name, tag);
//     });

//     tag.appendChild(text);
//     tag.appendChild(btn);

//     box.insertBefore(tag, input);

//     const hidden = document.createElement('input');
//     hidden.type = 'hidden';
//     hidden.name = 'ingredient_names';
//     hidden.value = name;
//     hidden.dataset.tag = name;
//     hiddenContainer.appendChild(hidden);
//   }

//   function removeIngredient(name, tagEl) {
//     tagEl.remove();
//     const hidden = document.querySelector(`#ingredient-hidden-inputs input[data-tag="${CSS.escape(name)}"]`);
//     if (hidden) hidden.remove();
//   }

  // document
  // .getElementById('form-product_create')
  // .addEventListener('submit', async (e) => {

  //   e.preventDefault();

  //   const name = document.getElementById('product-name').value;
  //   const origin = document.getElementById('product-origin').value;
  //   // Get ingredient list
  //   const ingredients = Array.from(
  //     document.querySelectorAll(
  //       '#ingredient-hidden-inputs input'
  //     )
  //   ).map(i => i.value);

  //   if (!name || ingredients.length == 0) {
  //     alert('Name and at least one ingredient required');
  //     return;
  //   }

  //   const ingredient_list = ingredients.join(', ');

    

  //     // QUERY PARAM VERSION
  //   const params = new URLSearchParams({
  //     name,
  //     origin,
  //     ingredient_list
  //   });

  //   const response = await fetch(
  //     `https://cff4ha.godiegogo.me/product/create?${params.toString()}`,
  //     {
  //       method: 'PUT'
  //     }
  //   );

  //   if (!response.ok) {
  //     throw new Error('Failed request');
  //   }

  //   const data = await response.json();
  //   console.log('SUCCESS:', data);
  //   alert('Product created successfully');

  //    //catch (err) {

  //   //   console.error(err);

  //   //   alert('Error creating product');
  //   // }
  // });

  document
  .getElementById('form-product_create')
  .addEventListener('submit', async (e) => {

    e.preventDefault();
    const name = document.getElementById('product-name').value;
    const origin = document.getElementById('product-origin').value;

    // Get ingredient list
    const ingredients = document.getElementById('ingredient-tag-input').value;
      

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

  // const ingredientInput = document.getElementById('ingredient-tag-input');
  // ingredientInput.addEventListener('input', handleIngredientInput);
  // ingredientInput.addEventListener('keydown', handleIngredientKeydown);

  document.getElementById('close-btn').addEventListener('click', () => {
    window.close();
  });