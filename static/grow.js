export function Growables() {
  document.querySelectorAll('.growable').forEach((e) => {
    let grow_event = "mouseenter";
    let shrink_event = "mouseleave";

    let target = e.parentNode;
    if (e.hasAttribute("data-growable-target")) {
      target = document.querySelector(e.getAttribute("data-growable-target"));
    }
    if (e.hasAttribute("data-growable-grow-event")) {
      grow_event = e.getAttribute("data-growable-grow-event");
    }
    if (e.hasAttribute("data-growable-shrink-event")) {
      shrink_event = e.getAttribute("data-growable-shrink-event");
    }

    target.addEventListener(grow_event, () => {
      e.classList.remove("growable");
      e.classList.add("growable--active");

      target.addEventListener(shrink_event, () => {
        e.classList.remove("growable--active");
        e.classList.add("growable");
      });
    });
  });
}

window.Growables = Growables;
document.addEventListener("htmx:load", Growables);
