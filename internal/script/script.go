package script

func GetInjectionScript() string {
	return `
(function(){
  var isMDeckHub = location.pathname === '/' || location.pathname === '/index.html' || location.href.includes('wails://') || location.href.includes('frontend/src');
  // --- 0. Zoom Manager (like Gotion) ---
  var zoomLevels = [0.7,0.8,0.85,0.9,0.95,1.0,1.05,1.1,1.2,1.35,1.5];
  var zoomIndex = 5;
  try{
    var saved = localStorage.getItem('mdeck_zoom_factor');
    if(saved){ var f=parseFloat(saved); for(var i=0;i<zoomLevels.length;i++){ if(Math.abs(zoomLevels[i]-f)<0.01){ zoomIndex=i; break; } } }
  }catch(e){}
  function showToast(msg){
    var t=document.getElementById('mdeck-toast');
    if(t) t.remove();
    var el=document.createElement('div');
    el.id='mdeck-toast';
    el.textContent=msg;
    el.style.cssText='position:fixed;bottom:22px;left:50%;transform:translateX(-50%);background:rgba(20,22,28,.92);color:#e8e8e5;padding:8px 14px;border-radius:8px;border:1px solid rgba(255,255,255,.08);font-size:12px;z-index:99999;backdrop-filter:blur(10px);pointer-events:none;font-family:Inter,system-ui,sans-serif;';
    document.body.appendChild(el);
    setTimeout(function(){ el.remove(); }, 1800);
  }
  function applyZoom(silent){
    var f=zoomLevels[zoomIndex];
    // target for zoom: hub grid or mdeck web content
    var target = document.getElementById('app') || document.documentElement;
    // Use CSS zoom (WebKit) with fallback transform
    if('zoom' in document.body.style){ document.documentElement.style.zoom = f; }
    else { document.documentElement.style.transform='scale('+f+')'; document.documentElement.style.transformOrigin='top left'; document.documentElement.style.width=(100/f)+'%'; document.documentElement.style.height=(100/f)+'%'; }
    // keep titlebar at 1x
    var bar=document.getElementById('mdeck-titlebar');
    if(bar){
      if('zoom' in bar.style) bar.style.zoom = 1/f;
      else { bar.style.transform='scale('+(1/f)+')'; bar.style.transformOrigin='top left'; bar.style.width=(100*f)+'%'; }
    }
    try{ localStorage.setItem('mdeck_zoom_factor', String(f)); }catch(e){}
    if(!silent) showToast('Zoom '+Math.round(f*100)+'%');
    // trigger xterm fit after zoom
    setTimeout(function(){ window.dispatchEvent(new Event('resize')); }, 80);
  }
  window.__mdeck_triggerZoomIn = function(){ if(zoomIndex<zoomLevels.length-1){ zoomIndex++; applyZoom(); } };
  window.__mdeck_triggerZoomOut = function(){ if(zoomIndex>0){ zoomIndex--; applyZoom(); } };
  window.__mdeck_triggerResetZoom = function(){ zoomIndex=5; applyZoom(); };
  setTimeout(function(){ applyZoom(true); }, 120);

  // --- 1. Titlebar + Resize Handles (seamless, like Gotion) ---
  if(document.getElementById('mdeck-titlebar')){ syncTheme(); return; }
  var styleId='mdeck-titlebar-style';
  if(!document.getElementById(styleId)){
    var s=document.createElement('style');
    s.id=styleId;
    s.textContent=[
      "#mdeck-titlebar{position:fixed;top:0;left:0;right:0;height:38px;display:flex;align-items:center;justify-content:space-between;padding:0 12px;background:rgba(12,14,18,.88);backdrop-filter:blur(16px) saturate(1.1);color:#e6e6e6;border-bottom:1px solid rgba(255,255,255,.06);z-index:2147483640;--wails-draggable:drag;user-select:none;font-family:Inter,system-ui,sans-serif;transition:background .18s,border-color .18s;}",
      "#mdeck-titlebar.light{background:rgba(246,245,244,.92);color:#31302e;border-bottom-color:#e6e6e6;}",
      ".mdeck-traffic{display:flex;gap:8px;--wails-draggable:no-drag}",
      ".mdeck-dot{width:12px;height:12px;border-radius:50%;border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:transform .12s,filter .15s;}",
      ".mdeck-dot:active{transform:scale(.9)}",
      ".mdeck-dot.close{background:#ff5f56;box-shadow:inset 0 0 0 1px #e0443e} .mdeck-dot.min{background:#ffbd2e;box-shadow:inset 0 0 0 1px #dea123} .mdeck-dot.max{background:#27c93f;box-shadow:inset 0 0 0 1px #1aab29}",
      ".mdeck-traffic .mdeck-dot svg{opacity:0;transition:opacity .15s} .mdeck-traffic:hover .mdeck-dot svg{opacity:.85}",
      ".mdeck-center{position:absolute;left:50%;transform:translateX(-50%);font-size:13px;font-weight:650;letter-spacing:-.2px;opacity:.92;pointer-events:none;display:flex;align-items:center;gap:8px}",
      ".mdeck-center small{font-weight:500;opacity:.55;font-size:11px}",
      ".mdeck-right{display:flex;gap:4px;--wails-draggable:no-drag;align-items:center}",
      ".mdeck-btn{width:28px;height:28px;border-radius:7px;border:none;background:transparent;color:#8b9099;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:background .12s,color .12s;}",
      ".mdeck-btn:hover{background:rgba(255,255,255,.08);color:#fff}",
      ".mdeck-btn:active{background:rgba(255,255,255,.12)}",
      ".mdeck-zoom{font-size:11px;color:#8b9099;padding:0 6px;min-width:44px;text-align:center;}",
      ".mdeck-resize{position:fixed;z-index:2147483647;background:transparent;--wails-draggable:no-drag}",
      ".mdeck-resize.top{top:0;left:10px;right:10px;height:5px;cursor:n-resize} .mdeck-resize.bottom{bottom:0;left:10px;right:10px;height:5px;cursor:s-resize}",
      ".mdeck-resize.left{top:10px;bottom:10px;left:0;width:5px;cursor:w-resize} .mdeck-resize.right{top:10px;bottom:10px;right:0;width:5px;cursor:e-resize}",
      ".mdeck-resize.tl{top:0;left:0;width:10px;height:10px;cursor:nw-resize} .mdeck-resize.tr{top:0;right:0;width:10px;height:10px;cursor:ne-resize}",
      ".mdeck-resize.bl{bottom:0;left:0;width:10px;height:10px;cursor:sw-resize} .mdeck-resize.br{bottom:0;right:0;width:10px;height:10px;cursor:se-resize}",
      "body.mdeck-has-titlebar{padding-top:38px !important;transition:padding-top .2s;}",
      "#app{transition:transform .12s;}",
      // polish for MDeck web tabs when inside wrapper
      ".mdeck-web-injected .ld-hud{top:38px !important;}",
    ].join("\n");
    (document.head||document.documentElement).appendChild(s);
  }
  function invoke(msg){
    try{ if(window.chrome&&window.chrome.webview&&window.chrome.webview.postMessage){ window.chrome.webview.postMessage(msg); return; } }catch(e){}
    try{ if(window.webkit&&window.webkit.messageHandlers&&window.webkit.messageHandlers.external){ window.webkit.messageHandlers.external.postMessage(msg); return; } }catch(e){}
    try{ if(window.WailsInvoke){ window.WailsInvoke(msg); return; } }catch(e){}
  }
  var bar=document.createElement('div'); bar.id='mdeck-titlebar';
  bar.innerHTML='<div class="mdeck-traffic">'
    +'<button class="mdeck-dot close" title="Close" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.Close&&window.go.main.App.Close()"><svg width=6 height=6 viewBox="0 0 6 6"><line x1=1 y1=1 x2=5 y2=5 stroke="#4d0000" stroke-width=1.2 stroke-linecap=round/><line x1=5 y1=1 x2=1 y2=5 stroke="#4d0000" stroke-width=1.2 stroke-linecap=round/></svg></button>'
    +'<button class="mdeck-dot min" title="Minimize" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.Minimise&&window.go.main.App.Minimise()"><svg width=6 height=6 viewBox="0 0 6 6"><line x1=1 y1=3 x2=5 y2=3 stroke="#664400" stroke-width=1.2 stroke-linecap=round/></svg></button>'
    +'<button class="mdeck-dot max" title="Maximize" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.ToggleMaximise&&window.go.main.App.ToggleMaximise()"><svg width=8 height=8 viewBox="0 0 12 12"><rect x=1.5 y=1.5 width=9 height=9 fill=none stroke="#0a4a0a" stroke-width=1.2/></svg></button>'
    +'</div><div class="mdeck-center"><span>MDeck</span><small id="mdeck-zoom-label">100%</small></div><div class="mdeck-right">'
    +'<button class="mdeck-btn" title="Zoom Out (Ctrl -)" onclick="window.__mdeck_triggerZoomOut&&window.__mdeck_triggerZoomOut()">−</button>'
    +'<span class="mdeck-zoom" id="mdeck-zoom-text">100%</span>'
    +'<button class="mdeck-btn" title="Zoom In (Ctrl +)" onclick="window.__mdeck_triggerZoomIn&&window.__mdeck_triggerZoomIn()">+</button>'
    +'<button class="mdeck-btn" title="Back (Alt ←)" onclick="history.back()">‹</button>'
    +'<button class="mdeck-btn" title="Reload (Ctrl R)" onclick="location.reload()">↻</button>'
    +'<button class="mdeck-btn" title="New Window (Ctrl N)" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.NewWindow&&window.go.main.App.NewWindow()">＋</button>'
    +'</div>';
  var last=0;
  bar.addEventListener('mousedown', function(e){
    if(e.target.closest('button')) return;
    if(e.button!==0) return;
    var now=Date.now(); if(now-last<280){ e.preventDefault(); if(window.go&&window.go.main&&window.go.main.App&&window.go.main.App.ToggleMaximise) window.go.main.App.ToggleMaximise(); return; } last=now;
    invoke('drag');
  });
  bar.addEventListener('dblclick', function(e){ if(e.target.closest('button')) return; if(window.go&&window.go.main&&window.go.main.App&&window.go.main.App.ToggleMaximise) window.go.main.App.ToggleMaximise(); });
  document.documentElement.appendChild(bar);
  document.body.classList.add('mdeck-has-titlebar');

  [['top'],['bottom'],['left'],['right'],['tl'],['tr'],['bl'],['br']].forEach(function(k){
    var cls=k[0]; if(document.querySelector('.mdeck-resize.'+cls)) return;
    var h=document.createElement('div'); h.className='mdeck-resize '+cls;
    h.addEventListener('mousedown', function(e){ if(e.button===0){ e.preventDefault(); var map={top:'n-resize',bottom:'s-resize',left:'w-resize',right:'e-resize',tl:'nw-resize',tr:'ne-resize',bl:'sw-resize',br:'se-resize'}; invoke('resize:'+map[cls]); }});
    document.documentElement.appendChild(h);
  });

  // --- 2. Global shortcuts (zoom, nav) ---
  if(!window.__mdeck_shortcuts){
    window.__mdeck_shortcuts=true;
    window.addEventListener('keydown', function(e){
      var mod=e.ctrlKey||e.metaKey;
      if(mod && (e.key==='='||e.key==='+'||e.code==='Equal'||e.code==='NumpadAdd')){ e.preventDefault(); window.__mdeck_triggerZoomIn(); return; }
      if(mod && (e.key==='-'||e.key==='_'||e.code==='Minus'||e.code==='NumpadSubtract')){ e.preventDefault(); window.__mdeck_triggerZoomOut(); return; }
      if(mod && (e.key==='0'||e.code==='Digit0'||e.code==='Numpad0')){ e.preventDefault(); window.__mdeck_triggerResetZoom(); return; }
      if(e.altKey && e.key==='ArrowLeft'){ e.preventDefault(); history.back(); }
      if(e.altKey && e.key==='ArrowRight'){ e.preventDefault(); history.forward(); }
      if(mod && e.key.toLowerCase()==='r' && !e.shiftKey){ e.preventDefault(); location.reload(); }
      if(mod && e.key.toLowerCase()==='n'){ e.preventDefault(); if(window.go&&window.go.main&&window.go.main.App&&window.go.main.App.NewWindow) window.go.main.App.NewWindow(); }
      // Ctrl +/- for xterm zoom is handled above; also handle Ctrl+Wheel
    }, true);
    window.addEventListener('wheel', function(e){
      if(e.ctrlKey){ e.preventDefault(); if(e.deltaY<0) window.__mdeck_triggerZoomIn(); else window.__mdeck_triggerZoomOut(); }
    }, {passive:false});
  }

  // keep zoom label in sync
  function syncZoomLabel(){
    var f=zoomLevels[zoomIndex];
    var el=document.getElementById('mdeck-zoom-text');
    var el2=document.getElementById('mdeck-zoom-label');
    if(el) el.textContent=Math.round(f*100)+'%';
    if(el2) el2.textContent=Math.round(f*100)+'%';
  }
  var origApply=applyZoom;
  applyZoom=function(silent){ origApply(silent); syncZoomLabel(); };
  syncZoomLabel();

  function syncTheme(){
    var isDark = !(document.body.classList.contains('light'));
    var b=document.getElementById('mdeck-titlebar');
    if(!b) return;
    // try to detect MDeck hub vs MDeck web
    var hub = !!document.getElementById('app');
    b.classList.toggle('light', !isDark);
    // also reflect zoom scale to keep 38px physical
  }
  syncTheme();
  setInterval(syncTheme, 900);
  // also observe window resize to keep seamless (like Gotion's titlebar counter-scale)
  if(!window.__mdeck_resize_obs){
    window.__mdeck_resize_obs=true;
    window.addEventListener('resize', function(){ syncTheme(); });
  }
  // mark web injected for MDeck tabs polish
  if(!isMDeckHub) document.documentElement.classList.add('mdeck-web-injected');
})();
`
}
