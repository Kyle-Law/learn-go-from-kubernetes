// Shared behaviour for lessons/NN-*/tutorial.html:
// syntax highlighting, copy buttons, table of contents, self-checks, progress bar.
(function () {
  const PAGE_KEY = "go-from-k8s-tutorial:" + location.pathname.split("/").slice(-2).join("/");

  function esc(s) {
    return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  }

  // ---------- Syntax highlighting ----------
  const GO_RE = new RegExp(
    [
      "(\\/\\/.*$)",                              // 1 comment
      "(\"(?:[^\"\\\\\\n]|\\\\.)*\"|`[^`]*`|'(?:[^'\\\\\\n]|\\\\.)')", // 2 string / raw string / rune
      "\\b(package|import|func|var|const|type|struct|interface|map|chan|for|range|if|else|switch|case|default|return|break|continue|go|defer|select|fallthrough|goto|nil|true|false|iota|make|new|append|len|cap|delete|panic|recover|any|string|bool|byte|rune|error|int|int32|int64|uint|float64)\\b", // 3 keyword
      "\\b(\\d+)\\b",                             // 4 number
    ].join("|"),
    "gm"
  );
  const HASH_RE = /(#.*$)|("(?:[^"\\\n]|\\.)*")/gm;

  function highlight(src, re) {
    let out = "", last = 0, m;
    re.lastIndex = 0;
    while ((m = re.exec(src))) {
      if (m[0] === "") { re.lastIndex++; continue; }
      out += esc(src.slice(last, m.index));
      const cls = m[1] ? "c" : m[2] ? "s" : m[3] ? "k" : "n";
      out += '<span class="tk-' + cls + '">' + esc(m[0]) + "</span>";
      last = re.lastIndex;
    }
    return out + esc(src.slice(last));
  }

  document.querySelectorAll("pre > code").forEach(code => {
    const pre = code.parentElement;
    const src = code.textContent.replace(/^\n/, "").replace(/\s+$/, "");
    if (pre.classList.contains("go")) code.innerHTML = highlight(src, GO_RE);
    else if (pre.classList.contains("sh") || pre.classList.contains("yaml")) code.innerHTML = highlight(src, HASH_RE);
    else code.textContent = src;

    // Wrap in .codeblock with optional label and a copy button.
    const wrap = document.createElement("div");
    wrap.className = "codeblock";
    pre.parentNode.insertBefore(wrap, pre);
    if (pre.dataset.label) {
      const label = document.createElement("span");
      label.className = "label";
      label.textContent = pre.dataset.label;
      wrap.appendChild(label);
    }
    wrap.appendChild(pre);
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "copy";
    btn.textContent = "Copy";
    btn.addEventListener("click", () => {
      const done = () => { btn.textContent = "Copied"; setTimeout(() => (btn.textContent = "Copy"), 1200); };
      if (navigator.clipboard) navigator.clipboard.writeText(src).then(done, () => {});
    });
    wrap.appendChild(btn);
  });

  // Inline `code` in .check option/explain text written with backticks.
  function inline(s) {
    return esc(s).replace(/`([^`]+)`/g, "<code>$1</code>");
  }

  // ---------- Table of contents ----------
  const tocList = document.querySelector(".toc ol");
  const tocDetails = document.querySelector(".toc details");
  // On narrow screens the contents sit above the article, so start collapsed.
  if (tocDetails && !matchMedia("(min-width: 64rem)").matches) tocDetails.open = false;
  const sections = Array.from(document.querySelectorAll("article section[id] > h2"));
  if (tocList) {
    sections.forEach(h => {
      const li = document.createElement("li");
      const a = document.createElement("a");
      a.href = "#" + h.parentElement.id;
      a.textContent = h.dataset.short || h.textContent;
      li.appendChild(a);
      tocList.appendChild(li);
    });
    const links = Array.from(tocList.querySelectorAll("a"));
    const setActive = () => {
      let current = sections[0];
      sections.forEach(h => { if (h.getBoundingClientRect().top < window.innerHeight * 0.3) current = h; });
      links.forEach(a => a.classList.toggle("active", current && a.getAttribute("href") === "#" + current.parentElement.id));
    };
    window.addEventListener("scroll", setActive, { passive: true });
    setActive();
  }

  // ---------- Progress bar ----------
  const bar = document.querySelector(".progress");
  if (bar) {
    const update = () => {
      const max = document.documentElement.scrollHeight - window.innerHeight;
      bar.style.width = (max > 0 ? (window.scrollY / max) * 100 : 0) + "%";
    };
    window.addEventListener("scroll", update, { passive: true });
    update();
  }

  // ---------- Self-checks ----------
  // Markup:
  // <div class="check" id="unique-id">
  //   <p class="prompt">Question</p>
  //   <pre class="go"><code>optional code</code></pre>
  //   <div class="options"><button data-ok>Right `answer`</button><button>Wrong</button></div>
  //   <div class="explain"><p>Why.</p></div>
  // </div>
  let saved = {};
  try { saved = JSON.parse(localStorage.getItem(PAGE_KEY)) || {}; } catch (e) { saved = {}; }
  const save = () => { try { localStorage.setItem(PAGE_KEY, JSON.stringify(saved)); } catch (e) { /* storage may be unavailable on file:// */ } };

  const checks = Array.from(document.querySelectorAll(".check[id]"));
  const scoreEl = document.querySelector(".toc .score");

  function updateScore() {
    if (!scoreEl) return;
    const answered = checks.filter(c => saved[c.id] !== undefined);
    const right = answered.filter(c => saved[c.id] === true).length;
    scoreEl.innerHTML = "Self-checks: <strong>" + right + "/" + checks.length + "</strong> correct" +
      (answered.length < checks.length ? " · " + (checks.length - answered.length) + " to go" : "");
  }

  checks.forEach(check => {
    const tag = document.createElement("span");
    tag.className = "tag";
    tag.textContent = "Check yourself";
    check.insertBefore(tag, check.firstChild);

    const buttons = Array.from(check.querySelectorAll(".options button"));
    buttons.forEach(b => { b.type = "button"; b.innerHTML = inline(b.textContent); });
    const explain = check.querySelector(".explain");
    const verdict = document.createElement("span");
    verdict.className = "verdict";
    const firstP = explain.querySelector("p");
    if (firstP) firstP.insertBefore(verdict, firstP.firstChild); else explain.prepend(verdict);
    const retry = document.createElement("button");
    retry.className = "retry";
    retry.type = "button";
    retry.textContent = "Try again";
    explain.appendChild(retry);

    function show(pickedIndex) {
      const ok = buttons[pickedIndex].hasAttribute("data-ok");
      buttons.forEach((b, i) => {
        b.disabled = true;
        b.classList.toggle("is-answer", b.hasAttribute("data-ok"));
        b.classList.toggle("is-picked-wrong", i === pickedIndex && !ok);
        b.classList.toggle("is-dim", !b.hasAttribute("data-ok") && i !== pickedIndex);
      });
      check.classList.toggle("correct", ok);
      check.classList.toggle("wrong", !ok);
      verdict.textContent = ok ? "Correct." : "Not quite.";
      explain.hidden = false;
      return ok;
    }
    function reset() {
      buttons.forEach(b => { b.disabled = false; b.className = ""; });
      check.classList.remove("correct", "wrong");
      explain.hidden = true;
    }

    explain.hidden = true;
    buttons.forEach((b, i) => b.addEventListener("click", () => {
      saved[check.id] = show(i);
      saved[check.id + ":i"] = i;
      save();
      updateScore();
    }));
    retry.addEventListener("click", () => {
      delete saved[check.id];
      delete saved[check.id + ":i"];
      save();
      reset();
      updateScore();
      buttons[0].focus();
    });

    const prev = saved[check.id + ":i"];
    if (Number.isInteger(prev) && buttons[prev]) show(prev);
  });
  updateScore();
})();
