export function Searchables() {
  document.querySelectorAll('.searchable').forEach((e) => {
    e.addEventListener("keyup", (event) => {
      let search_input = e;
      let search_empty_behavior = "hide-all";
      if (e.hasAttribute("data-searchable-empty-behavior")) {
        search_empty_behavior = e.getAttribute("data-searchable-empty-behavior");
      }

      if (search_input.value === "" || search_input.value === null) {
        let search_elements = e.parentNode.querySelectorAll(".searchable-item");
        if (search_empty_behavior === "show-all") {
          for (var i = 0; i < search_elements.length; i++) {
            search_elements[i].style.height = "";
            search_elements[i].style.width = "";
            search_elements[i].style.visibility = "visible";
          }
        } else if (search_empty_behavior === "hide-all") {
          for (var i = 0; i < search_elements.length; i++) {
            search_elements[i].style.height = "0px";
            search_elements[i].style.width = "0px";
            search_elements[i].style.margin = "0px";
            search_elements[i].style.padding = "0px";
            search_elements[i].style.opacity = "0";
            search_elements[i].style.visibility = "hidden";
          }
        }

        return;
      }

      let search_content = e.parentNode;
      if (e.hasAttribute("data-searchable-target")) {
        search_content = document.querySelector(e.getAttribute("data-searchable-target"));
      }

      // we are now going to iterate through all the children of search_content and hide those that 
      // don't match the search query, it's assumed that the searchable item has a data-searchable-data 
      // attribute that contains the value to search against.
      let search_elements = search_content.querySelectorAll(".searchable-item");
      for (var i = 0; i < search_elements.length; i++) {
        let element = search_elements[i];

        if (!element.hasAttribute("data-searchable-data") || element.getAttribute("data-searchable-data").toLowerCase().includes(search_input.value.toLowerCase()) === false) {
          // if it doesn't have the attribute, we will just hide it, 
          // as we have no way to search against it.
          element.style.height = "0px";
          element.style.width = "0px";
          element.style.margin = "0px";
          element.style.padding = "0px";
          element.style.opacity = "0";
          element.style.visibility = "hidden";
        } else if (element.getAttribute("data-searchable-data").toLowerCase().includes(search_input.value.toLowerCase())) {
          element.style.height = "";
          element.style.width = "";
          element.style.opacity = "1";

          element.style.visibility = "visible";
        }
      }
    });
  });
}

window.Searchables = Searchables;
