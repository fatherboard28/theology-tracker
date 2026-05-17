// app.js — minimal glue. Heavy lifting is done server-side via HTMX.

// ─── Sortable drag-and-drop for unit lists ───────────────────────────────────
// Any element with data-sortable will become a drag-and-drop list.
// On drop, HTMX posts the new order to data-sortable-url.
document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("[data-sortable]").forEach((el) => {
    Sortable.create(el, {
      handle: ".drag-handle",
      animation: 150,
      ghostClass: "sortable-ghost",
      onEnd(evt) {
        const url = el.dataset.sortableUrl;
        if (!url) return;

        const ids = [...el.querySelectorAll("[data-id]")].map(
          (item) => item.dataset.id
        );

        htmx.ajax("POST", url, {
          values: { ids: ids },
          target: el,
          swap: "none",
        });
      },
    });
  });
});

// ─── Markdown live preview ───────────────────────────────────────────────────
// Used by the note editor. Elements:
//   [data-md-source]  — the textarea
//   [data-md-preview] — the rendered output div
//
// Wire up once on load; HTMX re-wires after any hx-swap that recreates them.
function initMarkdownEditor(root = document) {
  const source = root.querySelector("[data-md-source]");
  const preview = root.querySelector("[data-md-preview]");
  if (!source || !preview) return;

  const render = () => {
    preview.innerHTML = marked.parse(source.value);
  };

  render();
  source.addEventListener("input", render);
}

document.addEventListener("DOMContentLoaded", () => initMarkdownEditor());
document.addEventListener("htmx:afterSwap", (e) => initMarkdownEditor(e.detail.elt));

// ─── Active nav link ─────────────────────────────────────────────────────────
function markActiveNav() {
  const path = window.location.pathname;
  document.querySelectorAll(".nav__links a").forEach((link) => {
    const href = link.getAttribute("href");
    const isActive = href === "/" ? path === "/" : path.startsWith(href);
    link.setAttribute("aria-current", isActive ? "page" : "false");
  });
}

document.addEventListener("DOMContentLoaded", markActiveNav);
document.addEventListener("htmx:afterSettle", markActiveNav);
