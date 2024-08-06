const refreshQuantity = (id) => {
  document.getElementById(`quantity-${id}`).innerHTML = "0";
  let product = document.getElementById(`product-${id}`);
  product.classList.remove("added");
};
