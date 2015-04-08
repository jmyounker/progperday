var ias = $('#scroller').ias({
  container:  ".container",
  item:       ".item",
  pagination: "#pagination",
  next:       ".next a",
  delay:      1250            // delay a little, so the user has
                              // more time to see the loader
});

// ias.extension(new IASSpinnerExtension());
// ias.extension(new IASNoneLeftExtension());

