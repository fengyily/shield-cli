// ══════════════════════════════════════════════
//  ER Diagram — Pure SVG, multi-layout, persistent positions
//  Adapted from PostgreSQL plugin for MySQL
// ══════════════════════════════════════════════

(function() {
  'use strict';

  const MIN_TABLE_WIDTH = 180;
  const MAX_TABLE_WIDTH = 400;
  const ROW_HEIGHT = 24;
  const HEADER_HEIGHT = 32;
  const PADDING = 40;
  const GAP_X = 50;
  const GAP_Y = 40;

  // State
  let erData = null;
  let erVisible = false;
  let viewX = 0, viewY = 0, scale = 1;
  let tablePositions = {};
  let tableWidths = {};
  let currentLayout = 'grid';

  // Interaction
  let panActive = false;
  let tableDrag = null;
  let colDrag = null;
  let selectedRel = null;
  let dragMoved = false;
  let hoveredCol = null;
  let dragStartX = 0, dragStartY = 0, dragViewX = 0, dragViewY = 0;
  let rafId = null;

  // Right-click pan
  let rightPanPending = null; // { sx, sy } screen coords at mousedown
  let rightPanned = false;    // true if right-drag exceeded threshold

  // Multi-select
  let selectedTables = new Set();  // set of table names
  let multiDragOffsets = {};       // { tableName: { dx, dy } } relative to drag anchor
  let marquee = null;              // { x1, y1, x2, y2 } in SVG coords when dragging
  let marqueePending = null;       // { x, y, sx, sy } waiting for threshold

  // ── Collaboration state ──
  let collabWs = null;
  let collabMe = null;
  let collabUsers = [];
  let collabCursors = {};
  let collabDrags = {};
  let collabReconnectTimer = null;
  let collabCursorCleanTimer = null;

  // ════════════════════════════════════════════
  //  Public API
  // ════════════════════════════════════════════

  window.toggleER = function() {
    const overlay = document.getElementById('erOverlay');
    erVisible = !erVisible;
    overlay.style.display = erVisible ? 'flex' : 'none';
    if (erVisible) {
      populateDbSelect();
      loadER();
      collabConnect();
      collabStartCleanup();
    } else {
      collabDisconnect();
      clearInterval(collabCursorCleanTimer);
    }
  };

  // Open ER diagram for a specific database (called from sidebar)
  window.openER = function(db) {
    const overlay = document.getElementById('erOverlay');
    erVisible = true;
    overlay.style.display = 'flex';
    populateDbSelect().then(() => {
      const sel = document.getElementById('erDbSelect');
      if (sel) sel.value = db;
      loadER();
      collabConnect();
      collabStartCleanup();
    });
  };

  window.loadER = async function() {
    const db = erCurrentDB();
    if (!db) return;
    const info = document.getElementById('erInfo');
    info.textContent = 'Loading...';

    try {
      const resp = await fetch('/api/er?db=' + encodeURIComponent(db));
      const data = await resp.json();
      if (data.code !== 200) { info.textContent = 'Error: ' + (data.message || 'unknown'); return; }

      erData = data.data;
      computeAllWidths();
      const tables = erData.tables || [];
      const rels = erData.relations || [];
      info.textContent = tables.length + ' tables, ' + rels.length + ' relations';

      if (!restorePositions()) {
        runLayout();
        fitToView();
      }
      applySvgTransform();
      renderSvg();
    } catch (e) {
      info.textContent = 'Failed to load';
    }
  };

  window.erSetLayout = function(mode) {
    currentLayout = mode;
    clearSavedPositions();
    computeAllWidths();
    runLayout();
    fitToView();
    renderSvg();
  };

  window.erClearLayout = function() {
    clearSavedPositions();
    document.getElementById('erLayoutSelect').value = currentLayout;
    computeAllWidths();
    runLayout();
    fitToView();
    renderSvg();
  };

  window.erZoom = function(factor) {
    const canvas = document.getElementById('erCanvas');
    const cx = canvas.clientWidth / 2, cy = canvas.clientHeight / 2;
    const ns = clamp(scale * factor);
    viewX = cx - (cx - viewX) * (ns / scale);
    viewY = cy - (cy - viewY) * (ns / scale);
    scale = ns;
    applySvgTransform();
  };

  window.erResetView = function() {
    if (erData) { fitToView(); applySvgTransform(); }
  };

  window.exportERSvg = function() {
    const svg = document.querySelector('#erCanvas svg');
    if (!svg) return;
    const clone = svg.cloneNode(true);
    const g = clone.querySelector('g');
    if (g) g.removeAttribute('transform');
    const bb = getContentBounds();
    const p = 30;
    clone.setAttribute('viewBox', (bb.x-p)+' '+(bb.y-p)+' '+(bb.w+p*2)+' '+(bb.h+p*2));
    clone.setAttribute('width', bb.w + p * 2);
    clone.setAttribute('height', bb.h + p * 2);
    clone.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    const blob = new Blob([clone.outerHTML], { type: 'image/svg+xml' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url; a.download = 'er-diagram.svg'; a.click();
    URL.revokeObjectURL(url);
  };

  // ════════════════════════════════════════════
  //  DB select (MySQL uses databases, not schemas)
  // ════════════════════════════════════════════

  function erCurrentDB() {
    const sel = document.getElementById('erDbSelect');
    return sel ? sel.value : (typeof currentDB !== 'undefined' ? currentDB : '');
  }

  async function populateDbSelect() {
    const sel = document.getElementById('erDbSelect');
    try {
      const resp = await fetch('/api/databases');
      const data = await resp.json();
      if (data.code === 200 && data.data) {
        sel.innerHTML = '';
        data.data.forEach(d => {
          const opt = document.createElement('option');
          opt.value = d; opt.textContent = d;
          if (d === (typeof currentDB !== 'undefined' ? currentDB : '')) opt.selected = true;
          sel.appendChild(opt);
        });
      }
    } catch(e) {}
  }

  // ════════════════════════════════════════════
  //  Persistence — localStorage
  // ════════════════════════════════════════════

  function storageKey() {
    const info = window._mysqlInfo || {};
    const db = erCurrentDB();
    return 'er:' + (info.host || '?') + ':' + (db || '?');
  }

  function savePositions() {
    try {
      const tables = (erData && erData.tables) || [];
      const tableList = tables.map(t => t.name).sort();
      const saved = {
        v: 1,
        layout: currentLayout,
        tables: {},
        tableList: tableList,
        view: { x: viewX, y: viewY, scale: scale },
        ts: Date.now()
      };
      tables.forEach(t => {
        const p = tablePositions[t.name];
        if (p) saved.tables[t.name] = { x: Math.round(p.x), y: Math.round(p.y) };
      });
      localStorage.setItem(storageKey(), JSON.stringify(saved));
    } catch(e) {}
  }

  function restorePositions() {
    try {
      const raw = localStorage.getItem(storageKey());
      if (!raw) return false;
      const saved = JSON.parse(raw);
      if (!saved || !saved.tables) return false;

      const tables = (erData && erData.tables) || [];
      if (tables.length === 0) return false;

      if (saved.layout) {
        currentLayout = saved.layout;
        document.getElementById('erLayoutSelect').value = currentLayout;
      }

      tablePositions = {};
      let restored = 0;
      tables.forEach(t => {
        if (saved.tables[t.name]) {
          tablePositions[t.name] = { x: saved.tables[t.name].x, y: saved.tables[t.name].y };
          restored++;
        }
      });
      if (restored === 0) return false;

      const newTables = tables.filter(t => !saved.tables[t.name]);
      if (newTables.length > 0) {
        placeNewTables(newTables);
      }

      if (saved.view) {
        viewX = saved.view.x || 0;
        viewY = saved.view.y || 0;
        scale = saved.view.scale || 1;
      }

      return true;
    } catch(e) { return false; }
  }

  function clearSavedPositions() {
    try { localStorage.removeItem(storageKey()); } catch(e) {}
  }

  function placeNewTables(newTables) {
    const bb = getContentBounds();
    const startY = bb.y + bb.h + GAP_Y * 2;
    const canvas = document.getElementById('erCanvas');
    const maxCols = Math.max(2, Math.floor(((canvas.clientWidth || 1200) - PADDING * 2 + GAP_X) / (avgTw() + GAP_X)));
    const lookup = {};
    (erData.tables || []).forEach(t => lookup[t.name] = t);
    placeGrid(newTables.map(t => t.name), maxCols, PADDING, startY, lookup);
  }

  // ════════════════════════════════════════════
  //  Layout dispatcher
  // ════════════════════════════════════════════

  function runLayout() {
    if (!erData || !(erData.tables || []).length) { tablePositions = {}; return; }
    tablePositions = {};

    switch (currentLayout) {
      case 'horizontal': layoutDirectional('h'); break;
      case 'vertical':   layoutDirectional('v'); break;
      case 'center':     layoutCenter(); break;
      default:           layoutGridMode(); break;
    }
  }

  function getTableHeight(t) {
    return HEADER_HEIGHT + (t.columns ? t.columns.length : 1) * ROW_HEIGHT;
  }

  const CHAR_W = 7;

  function computeTableWidth(t) {
    const headerW = t.name.length * 8 + 40;
    let maxColW = 0;
    (t.columns || []).forEach(c => {
      const nameW = c.name.length * CHAR_W + 30;
      const typeW = c.type.length * 6.5 + 10;
      maxColW = Math.max(maxColW, nameW + typeW + 24);
    });
    return Math.min(MAX_TABLE_WIDTH, Math.max(MIN_TABLE_WIDTH, headerW, maxColW));
  }

  function computeAllWidths() {
    if (!erData) return;
    (erData.tables || []).forEach(t => {
      tableWidths[t.name] = computeTableWidth(t);
    });
  }

  function tw(name) {
    return tableWidths[name] || MIN_TABLE_WIDTH;
  }

  function avgTw() {
    const vals = Object.values(tableWidths);
    return vals.length > 0 ? vals.reduce((a,b) => a+b, 0) / vals.length : MIN_TABLE_WIDTH;
  }

  function buildGraph() {
    const tables = erData.tables || [];
    const rels = erData.relations || [];
    const lookup = {};
    tables.forEach(t => lookup[t.name] = t);

    const adj = {};
    tables.forEach(t => adj[t.name] = new Set());
    rels.forEach(r => {
      if (adj[r.from_table]) adj[r.from_table].add(r.to_table);
      if (adj[r.to_table]) adj[r.to_table].add(r.from_table);
    });

    const fwd = {}, rev = {};
    tables.forEach(t => { fwd[t.name] = new Set(); rev[t.name] = new Set(); });
    rels.forEach(r => {
      if (fwd[r.from_table]) fwd[r.from_table].add(r.to_table);
      if (rev[r.to_table]) rev[r.to_table].add(r.from_table);
    });

    return { tables, rels, lookup, adj, fwd, rev };
  }

  function findComponents(tables, adj) {
    const visited = new Set();
    const multi = [], orphans = [];
    tables.forEach(t => {
      if (visited.has(t.name)) return;
      const comp = [];
      const q = [t.name]; visited.add(t.name);
      while (q.length) {
        const n = q.shift(); comp.push(n);
        adj[n].forEach(nb => { if (!visited.has(nb)) { visited.add(nb); q.push(nb); } });
      }
      if (comp.length > 1) multi.push(comp);
      else orphans.push(comp[0]);
    });
    multi.sort((a, b) => b.length - a.length);
    return { multi, orphans };
  }

  function canvasMaxCols() {
    const canvas = document.getElementById('erCanvas');
    return Math.max(2, Math.floor(((canvas.clientWidth || 1200) - PADDING * 2 + GAP_X) / (avgTw() + GAP_X)));
  }

  function placeGrid(names, cols, startX, startY, lookup) {
    let x = startX, y = startY, rowH = 0;
    names.forEach((name, i) => {
      const t = lookup[name];
      const h = t ? getTableHeight(t) : HEADER_HEIGHT + ROW_HEIGHT;
      tablePositions[name] = { x, y };
      rowH = Math.max(rowH, h);
      if ((i + 1) % cols === 0) { x = startX; y += rowH + GAP_Y; rowH = 0; }
      else x += tw(name) + GAP_X;
    });
    return y + rowH;
  }

  function layoutGridMode() {
    const { tables, lookup, adj } = buildGraph();
    const { multi, orphans } = findComponents(tables, adj);
    const maxCols = canvasMaxCols();
    let cy = PADDING;

    multi.forEach(comp => {
      const cols = Math.min(maxCols, Math.max(2, Math.ceil(Math.sqrt(comp.length * 1.5))));
      cy = placeGrid(comp, cols, PADDING, cy, lookup) + PADDING;
    });

    if (orphans.length > 0) {
      placeGrid(orphans, maxCols, PADDING, cy, lookup);
    }
  }

  function layoutDirectional(dir) {
    const g = buildGraph();
    const { tables, lookup, adj, fwd, rev } = g;
    const { multi, orphans } = findComponents(tables, adj);
    const maxCols = canvasMaxCols();

    let cursorMain = PADDING;

    multi.forEach(comp => {
      const layers = topoLayers(comp, fwd, rev);
      cursorMain = placeLayers(layers, dir, cursorMain, lookup);
      cursorMain += PADDING * 2;
    });

    if (orphans.length > 0) {
      placeGrid(orphans, maxCols, PADDING, cursorMain, lookup);
    }
  }

  function topoLayers(comp, fwd, rev) {
    const inComp = new Set(comp);
    const inDeg = {};
    comp.forEach(n => {
      inDeg[n] = 0;
      rev[n].forEach(p => { if (inComp.has(p)) inDeg[n]++; });
    });

    let queue = comp.filter(n => inDeg[n] === 0);
    if (queue.length === 0) queue = [comp[0]];

    const layers = [];
    const assigned = new Set();
    while (queue.length > 0) {
      layers.push([...queue]);
      queue.forEach(n => assigned.add(n));
      const next = [];
      queue.forEach(n => {
        fwd[n].forEach(child => {
          if (!inComp.has(child) || assigned.has(child)) return;
          inDeg[child]--;
          if (inDeg[child] <= 0 && !assigned.has(child)) next.push(child);
        });
      });
      queue = next;
    }
    const leftover = comp.filter(n => !assigned.has(n));
    if (leftover.length) layers.push(leftover);

    return layers;
  }

  function placeLayers(layers, dir, startOffset, lookup) {
    let mainCursor = startOffset;

    layers.forEach(layer => {
      let crossCursor = PADDING;
      let mainSize = 0;

      layer.forEach(name => {
        const t = lookup[name];
        const h = t ? getTableHeight(t) : HEADER_HEIGHT + ROW_HEIGHT;

        if (dir === 'h') {
          tablePositions[name] = { x: mainCursor, y: crossCursor };
          crossCursor += h + GAP_Y;
          mainSize = tw(name);
        } else {
          tablePositions[name] = { x: crossCursor, y: mainCursor };
          crossCursor += tw(name) + GAP_X;
          mainSize = Math.max(mainSize, h);
        }
      });

      mainCursor += (dir === 'h' ? avgTw() + GAP_X * 2 : mainSize + GAP_Y * 2);
    });

    return mainCursor;
  }

  function layoutCenter() {
    const g = buildGraph();
    const { tables, lookup, adj } = g;
    const { multi, orphans } = findComponents(tables, adj);

    let globalOffsetY = PADDING;

    multi.forEach(comp => {
      let centerName = comp[0], maxDeg = 0;
      comp.forEach(name => {
        const deg = adj[name].size;
        if (deg > maxDeg) { maxDeg = deg; centerName = name; }
      });

      const rings = [];
      const visited = new Set([centerName]);
      rings.push([centerName]);
      let frontier = [centerName];
      while (frontier.length) {
        const next = [];
        frontier.forEach(n => {
          adj[n].forEach(nb => {
            if (!visited.has(nb) && comp.includes(nb)) {
              visited.add(nb); next.push(nb);
            }
          });
        });
        if (next.length) rings.push(next);
        frontier = next;
      }

      const rBase = 0;
      const rStep = avgTw() * 1.6;
      const relPos = {};
      relPos[centerName] = { x: 0, y: 0 };

      for (let ri = 1; ri < rings.length; ri++) {
        const ring = rings[ri];
        const radius = rBase + ri * rStep;
        const angleStep = (2 * Math.PI) / ring.length;
        ring.forEach((name, idx) => {
          const angle = -Math.PI / 2 + idx * angleStep;
          relPos[name] = {
            x: Math.cos(angle) * radius,
            y: Math.sin(angle) * radius
          };
        });
      }

      let minRx = Infinity, minRy = Infinity, maxRx = -Infinity, maxRy = -Infinity;
      comp.forEach(name => {
        const p = relPos[name];
        if (!p) return;
        minRx = Math.min(minRx, p.x); minRy = Math.min(minRy, p.y);
        maxRx = Math.max(maxRx, p.x + tw(name)); maxRy = Math.max(maxRy, p.y);
      });

      const shiftX = PADDING - minRx;
      const shiftY = globalOffsetY - minRy;
      let compMaxY = 0;
      comp.forEach(name => {
        const p = relPos[name];
        if (!p) return;
        const t = lookup[name];
        const h = t ? getTableHeight(t) : HEADER_HEIGHT + ROW_HEIGHT;
        tablePositions[name] = { x: p.x + shiftX, y: p.y + shiftY };
        compMaxY = Math.max(compMaxY, p.y + shiftY + h);
      });

      globalOffsetY = compMaxY + PADDING * 2;
    });

    if (orphans.length > 0) {
      const maxCols = canvasMaxCols();
      placeGrid(orphans, maxCols, PADDING, globalOffsetY, lookup);
    }
  }

  // ════════════════════════════════════════════
  //  View helpers
  // ════════════════════════════════════════════

  function getContentBounds() {
    const tables = erData ? erData.tables || [] : [];
    let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
    tables.forEach(t => {
      const pos = tablePositions[t.name];
      if (!pos) return;
      minX = Math.min(minX, pos.x); minY = Math.min(minY, pos.y);
      maxX = Math.max(maxX, pos.x + tw(t.name));
      maxY = Math.max(maxY, pos.y + getTableHeight(t));
    });
    if (minX === Infinity) return { x: 0, y: 0, w: 0, h: 0 };
    return { x: minX, y: minY, w: maxX - minX, h: maxY - minY };
  }

  function fitToView() {
    const canvas = document.getElementById('erCanvas');
    const cw = canvas.clientWidth, ch = canvas.clientHeight;
    const bb = getContentBounds();
    if (bb.w === 0 || bb.h === 0) { scale = 1; viewX = 0; viewY = 0; return; }
    const pad = 40;
    const sx = (cw - pad * 2) / bb.w;
    const sy = (ch - pad * 2) / bb.h;
    scale = clamp(Math.min(sx, sy, 1.0));
    viewX = (cw - bb.w * scale) / 2 - bb.x * scale;
    viewY = (ch - bb.h * scale) / 2 - bb.y * scale;
  }

  function clamp(s) { return Math.max(0.08, Math.min(3, s)); }

  // ════════════════════════════════════════════
  //  Hit test
  // ════════════════════════════════════════════

  function screenToSvg(cx, cy) {
    const r = document.getElementById('erCanvas').getBoundingClientRect();
    return { x: (cx - r.left - viewX) / scale, y: (cy - r.top - viewY) / scale };
  }

  function hitTestTable(sx, sy) {
    if (!erData) return null;
    const tables = erData.tables || [];
    for (let i = tables.length - 1; i >= 0; i--) {
      const t = tables[i], pos = tablePositions[t.name];
      if (!pos) continue;
      if (sx >= pos.x && sx <= pos.x + tw(t.name) && sy >= pos.y && sy <= pos.y + getTableHeight(t))
        return t.name;
    }
    return null;
  }

  // ════════════════════════════════════════════
  //  SVG rendering
  // ════════════════════════════════════════════

  function renderSvg() {
    const canvas = document.getElementById('erCanvas');
    if (!erData) { canvas.innerHTML = ''; return; }
    const tables = erData.tables || [];
    const relations = erData.relations || [];
    if (tables.length === 0) { canvas.innerHTML = ''; return; }

    const tableLookup = {};
    tables.forEach(t => tableLookup[t.name] = t);

    let svg = '<svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg" style="position:absolute;top:0;left:0">';
    svg += svgDefs();
    svg += '<g transform="translate('+viewX+','+viewY+') scale('+scale+')">';
    svg += renderRelations(relations, tableLookup);
    tables.forEach(t => { svg += renderTableSvg(t); });
    if (colDrag) svg += renderDragLine(tableLookup);
    if (marquee) svg += renderMarquee();
    svg += renderCollabOverlays();
    svg += '</g></svg>';
    canvas.innerHTML = svg;
  }

  function renderDragLine(tableLookup) {
    if (!colDrag) return '';
    const fromPos = tablePositions[colDrag.fromTable];
    const fromTable = tableLookup[colDrag.fromTable];
    if (!fromPos || !fromTable) return '';
    const fi = fromTable.columns.findIndex(c => c.name === colDrag.fromCol);
    const y1 = fromPos.y + HEADER_HEIGHT + Math.max(0, fi) * ROW_HEIGHT + ROW_HEIGHT / 2;
    const x1 = fromPos.x + tw(colDrag.fromTable);

    let s = '<line x1="'+x1+'" y1="'+y1+'" x2="'+colDrag.curX+'" y2="'+colDrag.curY
      + '" stroke="#6366f1" stroke-width="2" stroke-dasharray="6,4" opacity="0.8"/>';

    const targetName = hitTestTable(colDrag.curX, colDrag.curY);
    if (targetName && targetName !== colDrag.fromTable) {
      const tp = tablePositions[targetName];
      const tt = tableLookup[targetName];
      if (tp && tt) {
        const th = getTableHeight(tt);
        const tW = tw(targetName);
        s += '<rect x="'+(tp.x-3)+'" y="'+(tp.y-3)+'" width="'+(tW+6)+'" height="'+(th+6)
          + '" rx="8" fill="none" stroke="#6366f1" stroke-width="2" stroke-dasharray="4,3" opacity="0.6"/>';

        const dropCol = hitTestColumn(colDrag.curX, colDrag.curY);
        if (dropCol && dropCol.table === targetName && !dropCol.pk) {
          const targetColObj = tt.columns.find(c => c.name === dropCol.column);
          const sourceTable = tableLookup[colDrag.fromTable];
          const sourceCol = sourceTable && sourceTable.columns.find(c => c.name === colDrag.fromCol);
          const typeMatch = targetColObj && sourceCol && targetColObj.type === sourceCol.type;
          const fillColor = typeMatch ? 'rgba(34,197,94,0.15)' : 'rgba(99,102,241,0.12)';
          s += '<rect x="'+tp.x+'" y="'+(tp.y + HEADER_HEIGHT + dropCol.idx * ROW_HEIGHT)
            + '" width="'+tW+'" height="'+ROW_HEIGHT+'" fill="'+fillColor+'" rx="3"/>';
        } else {
          const matchCol = tt.columns.find(c => c.name === colDrag.fromCol && !c.pk);
          if (matchCol) {
            const mi = tt.columns.findIndex(c => c.name === colDrag.fromCol);
            s += '<rect x="'+tp.x+'" y="'+(tp.y + HEADER_HEIGHT + mi * ROW_HEIGHT)
              + '" width="'+tW+'" height="'+ROW_HEIGHT+'" fill="rgba(99,102,241,0.12)" rx="3"/>';
          }
        }
      }
    }
    return s;
  }

  function renderMarquee() {
    if (!marquee) return '';
    const x = Math.min(marquee.x1, marquee.x2), y = Math.min(marquee.y1, marquee.y2);
    const w = Math.abs(marquee.x2 - marquee.x1), h = Math.abs(marquee.y2 - marquee.y1);
    return '<rect x="'+x+'" y="'+y+'" width="'+w+'" height="'+h
      + '" fill="rgba(76,110,245,0.08)" stroke="#4c6ef5" stroke-width="1" stroke-dasharray="6,3" rx="2"/>';
  }

  function svgDefs() {
    return '<defs>'
      + '<marker id="er-one" viewBox="0 0 10 12" refX="10" refY="6" markerWidth="8" markerHeight="10" orient="auto-start-reverse">'
      + '<line x1="8" y1="1" x2="8" y2="11" stroke="#7c8aa0" stroke-width="2"/></marker>'
      + '<marker id="er-many" viewBox="0 0 12 12" refX="0" refY="6" markerWidth="10" markerHeight="10" orient="auto-start-reverse">'
      + '<path d="M12,6 L2,1 M12,6 L2,11 M2,1 L2,11" fill="none" stroke="#7c8aa0" stroke-width="1.5"/></marker>'
      + '</defs>';
  }

  function renderRelations(relations, tableLookup) {
    let s = '';
    relations.forEach((rel, idx) => {
      const fp = tablePositions[rel.from_table], tp = tablePositions[rel.to_table];
      if (!fp || !tp) return;
      const ft = tableLookup[rel.from_table], tt = tableLookup[rel.to_table];
      if (!ft || !tt) return;

      const fi = ft.columns.findIndex(c => c.name === rel.from_column);
      const ti = tt.columns.findIndex(c => c.name === rel.to_column);
      const y1 = fp.y + HEADER_HEIGHT + Math.max(0, fi) * ROW_HEIGHT + ROW_HEIGHT / 2;
      const y2 = tp.y + HEADER_HEIGHT + Math.max(0, ti) * ROW_HEIGHT + ROW_HEIGHT / 2;

      const fw = tw(rel.from_table), tww = tw(rel.to_table);
      const fcx = fp.x + fw / 2, tcx = tp.x + tww / 2;
      let x1, x2;
      if (fcx <= tcx) { x1 = fp.x + fw; x2 = tp.x; }
      else { x1 = fp.x; x2 = tp.x + tww; }

      const dx = Math.max(Math.abs(x2 - x1) * 0.45, 30);
      const cx1 = x1 + (x1 < x2 ? dx : -dx), cx2 = x2 + (x1 < x2 ? -dx : dx);
      const path = 'M'+x1+','+y1+' C'+cx1+','+y1+' '+cx2+','+y2+' '+x2+','+y2;

      const isSel = selectedRel === idx;
      const stroke = isSel ? '#ef4444' : '#94a3b8';
      const width = isSel ? '2.5' : '1.5';
      const opacity = isSel ? '1' : '0.65';
      s += '<path d="'+path+'" fill="none" stroke="'+stroke+'" stroke-width="'+width+'" marker-start="url(#er-many)" marker-end="url(#er-one)" opacity="'+opacity+'"/>';
      s += '<path d="'+path+'" fill="none" stroke="transparent" stroke-width="14" data-rel-idx="'+idx+'" style="cursor:pointer"><title>'
        + escXml(rel.from_table+'.'+rel.from_column+' \u2192 '+rel.to_table+'.'+rel.to_column)+'</title></path>';
    });
    return s;
  }

  const ro = function() { return typeof readOnly !== 'undefined' && readOnly; };

  // MySQL header color: #00758f (MySQL blue/teal)
  function renderTableSvg(t) {
    const pos = tablePositions[t.name];
    if (!pos) return '';
    const W = tw(t.name);
    const cols = t.columns || [];
    const h = getTableHeight(t);
    let s = '';

    const isSel = selectedTables.has(t.name);
    s += '<rect x="'+(pos.x+3)+'" y="'+(pos.y+3)+'" width="'+W+'" height="'+h+'" rx="6" fill="rgba(0,0,0,0.06)"/>';
    if (isSel) {
      s += '<rect x="'+(pos.x-3)+'" y="'+(pos.y-3)+'" width="'+(W+6)+'" height="'+(h+6)+'" rx="8" fill="rgba(76,110,245,0.06)" stroke="#4c6ef5" stroke-width="2" stroke-dasharray="none"/>';
    }
    s += '<rect x="'+pos.x+'" y="'+pos.y+'" width="'+W+'" height="'+h+'" rx="6" fill="white" stroke="'+(isSel?'#4c6ef5':'#cbd5e1')+'" stroke-width="'+(isSel?'2':'1')+'"/>';
    s += '<path d="M'+(pos.x+6)+','+pos.y+' h'+(W-12)+' q6,0 6,6 v'+(HEADER_HEIGHT-6)+' h-'+W+' v-'+(HEADER_HEIGHT-6)+' q0,-6 6,-6 z" fill="#00758f"/>';
    s += '<line x1="'+pos.x+'" y1="'+(pos.y+HEADER_HEIGHT)+'" x2="'+(pos.x+W)+'" y2="'+(pos.y+HEADER_HEIGHT)+'" stroke="#cbd5e1" stroke-width="1"/>';
    s += '<text x="'+(pos.x+12)+'" y="'+(pos.y+21)+'" font-size="13" font-weight="600" fill="white" font-family="-apple-system,BlinkMacSystemFont,sans-serif">'
      + '\u25ED '+escXml(t.name)+'</text>';

    s += '<rect x="'+pos.x+'" y="'+pos.y+'" width="'+W+'" height="'+HEADER_HEIGHT+'" fill="transparent" style="cursor:move" data-table="'+escXml(t.name)+'"/>';

    if (!ro()) {
      s += '<g data-action="table-edit" data-table="'+escXml(t.name)+'" style="cursor:pointer">'
        + '<rect x="'+(pos.x+W-26)+'" y="'+(pos.y+2)+'" width="24" height="28" fill="transparent"/>'
        + '<text x="'+(pos.x+W-14)+'" y="'+(pos.y+22)+'" font-size="14" fill="white" text-anchor="middle" opacity="0.75" style="cursor:pointer">&#9881;</text>'
        + '</g>';
    }

    cols.forEach((col, i) => {
      const ry = pos.y + HEADER_HEIGHT + i * ROW_HEIGHT;
      if (i % 2 === 1) {
        if (i === cols.length - 1) {
          s += '<path d="M'+(pos.x+1)+','+ry+' h'+(W-2)+' v'+(ROW_HEIGHT-6)+' q0,5 -5,5 h-'+(W-12)+' q-5,0 -5,-5 z" fill="#f8f9fb"/>';
        } else {
          s += '<rect x="'+(pos.x+1)+'" y="'+ry+'" width="'+(W-2)+'" height="'+ROW_HEIGHT+'" fill="#f8f9fb"/>';
        }
      }
      const isPk = col.pk;
      if (isPk) {
        s += '<text x="'+(pos.x+10)+'" y="'+(ry+16)+'" font-size="11" fill="#e67700">\uD83D\uDD11</text>';
      }
      const isFk = erData.relations && erData.relations.some(r => r.from_table === t.name && r.from_column === col.name);
      if (isFk && !isPk) {
        s += '<text x="'+(pos.x+10)+'" y="'+(ry+16)+'" font-size="10" fill="#6366f1">\u2192</text>';
      }
      const nx = pos.x + ((isPk || isFk) ? 26 : 12);
      s += '<text x="'+nx+'" y="'+(ry+16)+'" font-size="12" fill="#1e293b" font-family="-apple-system,BlinkMacSystemFont,sans-serif" font-weight="'+(isPk?'600':'400')+'">'
        + escXml(col.name)+'</text>';

      s += '<text x="'+(pos.x+W-10)+'" y="'+(ry+16)+'" font-size="11" fill="#94a3b8" text-anchor="end" font-family="-apple-system,BlinkMacSystemFont,sans-serif">'
        + escXml(col.type)+'</text>';

      if (!ro() && hoveredCol && hoveredCol.table === t.name && hoveredCol.idx === i) {
        s += '<g data-action="col-edit" data-table="'+escXml(t.name)+'" data-col="'+i+'" style="cursor:pointer">'
          + '<rect x="'+(pos.x+W-20)+'" y="'+(ry+1)+'" width="18" height="22" fill="white" opacity="0.85" rx="3"/>'
          + '<text x="'+(pos.x+W-11)+'" y="'+(ry+17)+'" font-size="13" fill="#64748b" text-anchor="middle" style="cursor:pointer">&#9881;</text>'
          + '</g>';
      }
    });

    return s;
  }

  function applySvgTransform() {
    const g = document.querySelector('#erCanvas svg > g');
    if (g) g.setAttribute('transform', 'translate('+viewX+','+viewY+') scale('+scale+')');
  }

  function escXml(s) {
    if (s == null) return '';
    return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
  }

  // ════════════════════════════════════════════
  //  Interaction: Pan + Table drag + Column drag + Zoom
  // ════════════════════════════════════════════

  const canvas = document.getElementById('erCanvas');

  function scheduleRender() {
    if (rafId) return;
    rafId = requestAnimationFrame(() => { rafId = null; renderSvg(); });
  }

  function hitTestColumn(sx, sy) {
    if (!erData) return null;
    const tables = erData.tables || [];
    for (let i = tables.length - 1; i >= 0; i--) {
      const t = tables[i], pos = tablePositions[t.name];
      if (!pos) continue;
      if (sx < pos.x || sx > pos.x + tw(t.name)) continue;
      const colY = sy - pos.y - HEADER_HEIGHT;
      if (colY < 0) continue;
      const colIdx = Math.floor(colY / ROW_HEIGHT);
      if (colIdx >= 0 && colIdx < (t.columns || []).length) {
        const col = t.columns[colIdx];
        return { table: t.name, column: col.name, pk: col.pk, idx: colIdx };
      }
    }
    return null;
  }

  let pendingIcon = null;

  canvas.addEventListener('mousedown', e => {
    // Middle mouse button or right mouse button → pan (right uses threshold)
    if (e.button === 1) {
      e.preventDefault();
      panActive = true;
      dragStartX = e.clientX; dragStartY = e.clientY;
      dragViewX = viewX; dragViewY = viewY;
      canvas.style.cursor = 'grabbing';
      return;
    }
    if (e.button === 2) {
      rightPanPending = { sx: e.clientX, sy: e.clientY };
      rightPanned = false;
      dragViewX = viewX; dragViewY = viewY;
      canvas.style.cursor = 'grab';
      return;
    }
    if (e.button !== 0) return;
    dragMoved = false;
    pendingIcon = null;
    const pt = screenToSvg(e.clientX, e.clientY);

    const g = e.target.closest('[data-action]');
    if (g) {
      pendingIcon = {
        action: g.getAttribute('data-action'),
        table: g.getAttribute('data-table'),
        col: g.hasAttribute('data-col') ? parseInt(g.getAttribute('data-col'), 10) : -1,
        clientX: e.clientX, clientY: e.clientY
      };
      e.preventDefault();
      return;
    }

    e.preventDefault();

    // Skip selection logic when clicking on a relation line
    if (e.target.closest('[data-rel-idx]')) return;

    const tableHit = hitTestTable(pt.x, pt.y);
    const colHit = hitTestColumn(pt.x, pt.y);

    if (tableHit) {
      const tpos = tablePositions[tableHit];

      // Shift/Cmd+Click toggles multi-selection (anywhere on table)
      if (e.shiftKey || e.metaKey) {
        if (selectedTables.has(tableHit)) {
          selectedTables.delete(tableHit);
        } else {
          selectedTables.add(tableHit);
        }
        scheduleRender();
        return;
      }

      // Single click — select this table
      if (!selectedTables.has(tableHit)) {
        selectedTables.clear();
        selectedTables.add(tableHit);
        scheduleRender();
      }

      // Column area → FK drag (always, regardless of selection)
      if (colHit && !ro()) {
        colDrag = { fromTable: colHit.table, fromCol: colHit.column, curX: pt.x, curY: pt.y };
        canvas.style.cursor = 'crosshair';
        return;
      }

      // Header area → table drag (moves all selected)
      if (pt.y < tpos.y + HEADER_HEIGHT) {
        multiDragOffsets = {};
        selectedTables.forEach(name => {
          const p = tablePositions[name];
          if (p) {
            multiDragOffsets[name] = { dx: p.x - tpos.x, dy: p.y - tpos.y };
          }
        });
        tableDrag = { name: tableHit, ox: pt.x - tpos.x, oy: pt.y - tpos.y };
        canvas.style.cursor = 'grabbing';
      }
      return;
    }

    // Empty space — prepare marquee (activated after threshold in mousemove)
    if (!e.shiftKey && !e.metaKey) {
      selectedTables.clear();
      scheduleRender();
    }
    marqueePending = { x: pt.x, y: pt.y, sx: e.clientX, sy: e.clientY };
  });

  window.addEventListener('mousemove', e => {
    if (colDrag) {
      dragMoved = true;
      const pt = screenToSvg(e.clientX, e.clientY);
      colDrag.curX = pt.x; colDrag.curY = pt.y;
      collabSendCursor(pt.x, pt.y);
      scheduleRender();
    } else if (tableDrag) {
      dragMoved = true;
      const pt = screenToSvg(e.clientX, e.clientY);
      const newPos = { x: pt.x - tableDrag.ox, y: pt.y - tableDrag.oy };
      tablePositions[tableDrag.name] = newPos;
      collabSendDrag(tableDrag.name, newPos.x, newPos.y, false);
      // Move other selected tables with the same delta
      if (selectedTables.size > 0) {
        selectedTables.forEach(name => {
          if (name !== tableDrag.name && multiDragOffsets[name]) {
            const off = multiDragOffsets[name];
            const p = { x: newPos.x + off.dx, y: newPos.y + off.dy };
            tablePositions[name] = p;
            collabSendDrag(name, p.x, p.y, false);
          }
        });
      }
      scheduleRender();
    } else if (marqueePending && !marquee) {
      // Activate marquee only after 4px movement threshold
      const dx = e.clientX - marqueePending.sx, dy = e.clientY - marqueePending.sy;
      if (dx * dx + dy * dy >= 16) {
        marquee = { x1: marqueePending.x, y1: marqueePending.y, x2: marqueePending.x, y2: marqueePending.y };
        marqueePending = null;
        canvas.style.cursor = 'crosshair';
      }
    } else if (marquee) {
      dragMoved = true;
      const pt = screenToSvg(e.clientX, e.clientY);
      marquee.x2 = pt.x; marquee.y2 = pt.y;
      const mx1 = Math.min(marquee.x1, marquee.x2), my1 = Math.min(marquee.y1, marquee.y2);
      const mx2 = Math.max(marquee.x1, marquee.x2), my2 = Math.max(marquee.y1, marquee.y2);
      selectedTables.clear();
      (erData.tables || []).forEach(t => {
        const p = tablePositions[t.name];
        if (!p) return;
        const tx2 = p.x + tw(t.name), ty2 = p.y + getTableHeight(t);
        if (p.x < mx2 && tx2 > mx1 && p.y < my2 && ty2 > my1) selectedTables.add(t.name);
      });
      scheduleRender();
    } else if (panActive) {
      dragMoved = true;
      viewX = dragViewX + (e.clientX - dragStartX);
      viewY = dragViewY + (e.clientY - dragStartY);
      applySvgTransform();
    } else if (rightPanPending) {
      const dx = e.clientX - rightPanPending.sx, dy = e.clientY - rightPanPending.sy;
      if (!rightPanned && dx * dx + dy * dy >= 16) {
        rightPanned = true;
        dragStartX = rightPanPending.sx; dragStartY = rightPanPending.sy;
        canvas.style.cursor = 'grabbing';
      }
      if (rightPanned) {
        viewX = dragViewX + (e.clientX - dragStartX);
        viewY = dragViewY + (e.clientY - dragStartY);
        applySvgTransform();
      }
    }
  });

  window.addEventListener('mouseup', () => {
    if (pendingIcon && !dragMoved && !ro()) {
      const info = pendingIcon;
      pendingIcon = null;
      const tbl = (erData.tables || []).find(t => t.name === info.table);
      if (tbl && info.action === 'table-edit') {
        showTableStructureDialog(info.table);
      } else if (tbl && info.action === 'col-edit' && info.col >= 0 && tbl.columns[info.col]) {
        const col = tbl.columns[info.col];
        showCtxMenu(info.clientX, info.clientY, [
          { label: 'Edit Column...', icon: '&#9998;', action: () => showEditColumnDialog(info.table, col) },
          { label: 'Delete Column', icon: '&#128465;', danger: true, action: () => showDeleteColumnDialog(info.table, col) },
          '---',
          { label: 'Add Column...', icon: '&#10010;', action: () => showAddColumnDialog(info.table) },
        ]);
      }
      return;
    }
    pendingIcon = null;

    if (colDrag) {
      const pt = { x: colDrag.curX, y: colDrag.curY };
      const targetName = hitTestTable(pt.x, pt.y);
      if (targetName && targetName !== colDrag.fromTable) {
        const targetTable = (erData.tables || []).find(t => t.name === targetName);
        const sourceTable = (erData.tables || []).find(t => t.name === colDrag.fromTable);
        const sourceCol = sourceTable && sourceTable.columns.find(c => c.name === colDrag.fromCol);
        if (targetTable && sourceCol) {
          const dropCol = hitTestColumn(pt.x, pt.y);
          if (dropCol && dropCol.table === targetName && !dropCol.pk) {
            const targetColObj = targetTable.columns.find(c => c.name === dropCol.column);
            if (targetColObj && targetColObj.type === sourceCol.type) {
              showFkConfirmDialog({
                fromTable: colDrag.fromTable, fromCol: colDrag.fromCol,
                toTable: targetName, toCol: dropCol.column,
                colType: sourceCol.type, needsCreate: false
              });
            } else {
              const defaultName = colDrag.fromTable + '_' + colDrag.fromCol;
              showFkConfirmDialog({
                fromTable: colDrag.fromTable, fromCol: colDrag.fromCol,
                toTable: targetName, toCol: defaultName,
                colType: sourceCol.type, needsCreate: true
              });
            }
          } else {
            const matchCol = targetTable.columns.find(c => c.name === colDrag.fromCol && !c.pk);
            if (matchCol) {
              showFkConfirmDialog({
                fromTable: colDrag.fromTable, fromCol: colDrag.fromCol,
                toTable: targetName, toCol: colDrag.fromCol,
                colType: sourceCol.type, needsCreate: false
              });
            } else {
              const defaultName = colDrag.fromTable + '_' + colDrag.fromCol;
              showFkConfirmDialog({
                fromTable: colDrag.fromTable, fromCol: colDrag.fromCol,
                toTable: targetName, toCol: defaultName,
                colType: sourceCol.type, needsCreate: true
              });
            }
          }
        }
      }
      colDrag = null;
      canvas.style.cursor = 'default';
      scheduleRender();
    } else if (tableDrag) {
      if (dragMoved) {
        const pos = tablePositions[tableDrag.name];
        if (pos) collabSendDrag(tableDrag.name, pos.x, pos.y, true);
        // Send done for other selected tables too
        if (selectedTables.size > 0) {
          selectedTables.forEach(name => {
            if (name !== tableDrag.name) {
              const p = tablePositions[name];
              if (p) collabSendDrag(name, p.x, p.y, true);
            }
          });
        }
        savePositions();
      }
      multiDragOffsets = {};
      tableDrag = null;
      canvas.style.cursor = 'default';
    } else if (marquee || marqueePending) {
      marquee = null;
      marqueePending = null;
      canvas.style.cursor = 'default';
      scheduleRender();
    } else if (panActive) {
      savePositions();
      panActive = false;
      canvas.style.cursor = 'default';
    } else if (rightPanPending) {
      const wasDrag = rightPanned;
      const sx = rightPanPending.sx, sy = rightPanPending.sy;
      if (wasDrag) savePositions();
      rightPanPending = null;
      rightPanned = false;
      canvas.style.cursor = 'default';
      // Right-click without drag → show context menu
      if (!wasDrag) showErContextMenu(sx, sy);
    }
  });

  canvas.addEventListener('click', e => {
    const el = e.target.closest('[data-rel-idx]');
    if (el) {
      const idx = parseInt(el.getAttribute('data-rel-idx'), 10);
      selectedRel = (selectedRel === idx) ? null : idx;
      scheduleRender();
    } else if (selectedRel !== null) {
      const pt = screenToSvg(e.clientX, e.clientY);
      const tbl = hitTestTable(pt.x, pt.y);
      if (!tbl) {
        selectedRel = null;
        scheduleRender();
      }
    }
  });

  document.addEventListener('keydown', e => {
    if (!erVisible) return;
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.tagName === 'SELECT') return;

    // Delete/Backspace → delete selected relation
    if (selectedRel !== null && !ro()) {
      if (e.key === 'Delete' || e.key === 'Backspace') {
        e.preventDefault();
        const rels = (erData && erData.relations) || [];
        const rel = rels[selectedRel];
        if (rel) showFkDeleteDialog(rel);
      }
    }
  });

  canvas.addEventListener('wheel', e => {
    e.preventDefault();
    if (e.ctrlKey || e.metaKey) {
      const r = canvas.getBoundingClientRect();
      const mx = e.clientX - r.left, my = e.clientY - r.top;
      const f = e.deltaY < 0 ? 1.12 : 1 / 1.12;
      const ns = clamp(scale * f);
      viewX = mx - (mx - viewX) * (ns / scale);
      viewY = my - (my - viewY) * (ns / scale);
      scale = ns;
    } else {
      viewX -= e.deltaX;
      viewY -= e.deltaY;
    }
    applySvgTransform();
    debounceSave();
  }, { passive: false });

  let saveTimer = null;
  function debounceSave() {
    clearTimeout(saveTimer);
    saveTimer = setTimeout(savePositions, 300);
  }

  canvas.addEventListener('mousemove', function hoverCursor(e) {
    if (tableDrag || marquee || colDrag || rightPanPending) return;
    const pt = screenToSvg(e.clientX, e.clientY);
    collabSendCursor(pt.x, pt.y);
    const col = hitTestColumn(pt.x, pt.y);

    const prevHover = hoveredCol;
    if (col) {
      hoveredCol = { table: col.table, idx: col.idx };
      const iconG = e.target.closest('[data-action]');
      canvas.style.cursor = iconG ? 'pointer' : 'crosshair';
    } else {
      hoveredCol = null;
      if (hitTestTable(pt.x, pt.y)) canvas.style.cursor = 'move';
      else canvas.style.cursor = 'default';
    }

    if (!ro()) {
      const changed = (!prevHover && hoveredCol) || (prevHover && !hoveredCol)
        || (prevHover && hoveredCol && (prevHover.table !== hoveredCol.table || prevHover.idx !== hoveredCol.idx));
      if (changed) scheduleRender();
    }
  });

  document.addEventListener('keydown', e => {
    if (e.key === 'Escape' && erVisible) toggleER();
  });

  // ── Toast ──

  function showToast(msg, type) {
    const t = document.getElementById('toast');
    if (!t) { alert(msg); return; }
    t.textContent = msg;
    t.className = 'toast show' + (type ? ' toast-' + type : '');
    clearTimeout(t._timer);
    t._timer = setTimeout(() => { t.className = 'toast'; }, 3000);
  }

  // ── MySQL identifier quoting (backticks) ──

  function quoteId(s) { return '`' + String(s).replace(/`/g, '``') + '`'; }

  function erDB() { return erCurrentDB(); }

  // ── Shared SQL executor ──

  async function execSql(sql, db) {
    const resp = await fetch('/api/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ sql: sql, db: db || erDB() })
    });
    const data = await resp.json();
    if (data.code !== 200) throw new Error(data.message || 'unknown error');
    return data;
  }

  async function execStmts(stmts, db) {
    for (const sql of stmts) await execSql(sql, db);
  }

  // ── Common MySQL type list ──

  const MYSQL_TYPES = [
    'int','bigint','smallint','tinyint','mediumint',
    'decimal','float','double',
    'varchar(255)','char(50)','text','mediumtext','longtext','tinytext',
    'boolean',
    'datetime','timestamp','date','time','year',
    'json','blob','mediumblob','longblob',
    'enum','set','binary','varbinary(255)'
  ];

  function mysqlTypeSelect(selected, id) {
    let s = '<select id="'+id+'" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none">';
    MYSQL_TYPES.forEach(t => {
      s += '<option value="'+escXml(t)+'"'+(t===selected?' selected':'')+'>'+escXml(t)+'</option>';
    });
    s += '</select>';
    return s;
  }

  // ── FK creation dialog ──

  function showFkConfirmDialog(opts) {
    const db = erDB();

    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.style.zIndex = '400';

    function buildSql(fkColName) {
      const cname = 'fk_' + opts.toTable + '_' + fkColName + '_' + opts.fromTable;
      const stmts = [];
      if (opts.needsCreate) {
        stmts.push('ALTER TABLE ' + quoteId(opts.toTable)
          + ' ADD COLUMN ' + quoteId(fkColName) + ' ' + opts.colType);
      }
      stmts.push('ALTER TABLE ' + quoteId(opts.toTable)
        + ' ADD CONSTRAINT ' + quoteId(cname)
        + ' FOREIGN KEY (' + quoteId(fkColName) + ') REFERENCES '
        + quoteId(opts.fromTable) + '(' + quoteId(opts.fromCol) + ')');
      return stmts;
    }

    const diagramHtml = '<div style="display:flex;align-items:center;gap:10px;padding:14px;background:#f8f9fb;border-radius:6px;border:1px solid #e0e2e8;margin-bottom:16px;font-size:13px">'
      + '<div style="font-family:monospace"><strong>' + escXml(opts.fromTable) + '</strong>.' + escXml(opts.fromCol)
      + ' <span style="color:#94a3b8">(' + escXml(opts.colType) + ')</span></div>'
      + '<div style="color:#6366f1;font-size:18px">&larr;</div>'
      + '<div style="font-family:monospace"><strong>' + escXml(opts.toTable) + '</strong>.<span id="erFkPreviewCol">' + escXml(opts.toCol) + '</span></div>'
      + '</div>';

    const colInputHtml = opts.needsCreate
      ? '<div style="margin-bottom:14px">'
        + '<label style="display:block;font-size:12px;color:#6b7084;margin-bottom:5px">FK column name in <strong>' + escXml(opts.toTable) + '</strong> (new column will be created)</label>'
        + '<input type="text" id="erFkColName" value="' + escXml(opts.toCol) + '" '
        + 'style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/>'
        + '</div>'
      : '<div style="margin-bottom:14px;font-size:12px;color:#6b7084">'
        + 'Using existing column <strong style="font-family:monospace">' + escXml(opts.toCol) + '</strong> in ' + escXml(opts.toTable)
        + '</div>';

    overlay.innerHTML = '<div class="modal" style="width:560px">'
      + '<div class="modal-header"><span>Create Foreign Key</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + diagramHtml + colInputHtml
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erFkSqlPre" style="margin-top:8px;padding:10px;background:#f8f9fb;border-radius:4px;font-size:11px;overflow-x:auto;border:1px solid #e0e2e8;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erFkConfirmBtn">Create</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    function updateSqlPreview() {
      const colName = opts.needsCreate
        ? (document.getElementById('erFkColName').value || opts.toCol)
        : opts.toCol;
      const stmts = buildSql(colName);
      document.getElementById('erFkSqlPre').textContent = stmts.join(';\n') + ';';
      const preview = document.getElementById('erFkPreviewCol');
      if (preview) preview.textContent = colName;
    }
    updateSqlPreview();

    if (opts.needsCreate) {
      const input = document.getElementById('erFkColName');
      input.addEventListener('input', updateSqlPreview);
      input.focus();
      input.select();
    }

    document.getElementById('erFkConfirmBtn').onclick = async function() {
      const colName = opts.needsCreate
        ? (document.getElementById('erFkColName').value || '').trim()
        : opts.toCol;

      if (!colName) {
        showToast('Column name cannot be empty', 'warning');
        return;
      }

      this.disabled = true; this.textContent = 'Creating...';

      const stmts = buildSql(colName);
      let success = true;
      for (const sql of stmts) {
        try {
          const resp = await fetch('/api/query', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ sql: sql, db: db })
          });
          const data = await resp.json();
          if (data.code !== 200) {
            showToast('Error: ' + (data.message || 'unknown'), 'error');
            this.disabled = false; this.textContent = 'Create';
            success = false;
            break;
          }
        } catch(e) {
          showToast('Request failed', 'error');
          this.disabled = false; this.textContent = 'Create';
          success = false;
          break;
        }
      }

      if (success) {
        overlay.remove();
        showToast('Foreign key created' + (opts.needsCreate ? ' (new column added)' : ''), 'success');
        collabNotifySchemaChanged();
        loadER();
      }
    };
  }

  // ── FK delete dialog (MySQL uses DROP FOREIGN KEY, not DROP CONSTRAINT) ──

  function showFkDeleteDialog(rel) {
    const db = erDB();
    const constraintName = rel.constraint || '';

    if (!constraintName) {
      showToast('Cannot determine constraint name', 'error');
      return;
    }

    const sql = 'ALTER TABLE ' + quoteId(rel.from_table) + ' DROP FOREIGN KEY ' + quoteId(constraintName);

    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.style.zIndex = '400';

    const diagramHtml = '<div style="display:flex;align-items:center;gap:10px;padding:14px;background:#fef2f2;border-radius:6px;border:1px solid #fecaca;margin-bottom:16px;font-size:13px">'
      + '<div style="font-family:monospace"><strong>' + escXml(rel.from_table) + '</strong>.' + escXml(rel.from_column) + '</div>'
      + '<div style="color:#ef4444;font-size:18px">&rarr;</div>'
      + '<div style="font-family:monospace"><strong>' + escXml(rel.to_table) + '</strong>.' + escXml(rel.to_column) + '</div>'
      + '</div>';

    overlay.innerHTML = '<div class="modal" style="width:520px">'
      + '<div class="modal-header"><span>Drop Foreign Key</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<p style="margin:0 0 12px;color:#64748b;font-size:13px">Constraint: <strong style="font-family:monospace;color:#1e293b">' + escXml(constraintName) + '</strong></p>'
      + diagramHtml
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre style="margin-top:8px;padding:10px;background:#fef2f2;border-radius:4px;font-size:11px;overflow-x:auto;border:1px solid #fecaca;white-space:pre-wrap">' + escXml(sql) + ';</pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erFkDeleteBtn" style="background:#ef4444;border-color:#ef4444">Drop</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    document.getElementById('erFkDeleteBtn').onclick = async function() {
      this.disabled = true; this.textContent = 'Dropping...';
      try {
        const resp = await fetch('/api/query', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ sql: sql, db: db })
        });
        const data = await resp.json();
        if (data.code !== 200) {
          showToast('Error: ' + (data.message || 'unknown'), 'error');
          this.disabled = false; this.textContent = 'Drop';
          return;
        }
        overlay.remove();
        selectedRel = null;
        showToast('Foreign key dropped: ' + constraintName, 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Request failed', 'error');
        this.disabled = false; this.textContent = 'Drop';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Context Menu
  // ══════════════════════════════════════════════

  let ctxMenu = null;

  function hideCtxMenu() {
    if (ctxMenu) { ctxMenu.remove(); ctxMenu = null; }
  }

  function showCtxMenu(x, y, items) {
    hideCtxMenu();
    const div = document.createElement('div');
    div.className = 'er-ctx';
    items.forEach(item => {
      if (item === '---') {
        const sep = document.createElement('div');
        sep.className = 'er-ctx-sep';
        div.appendChild(sep);
        return;
      }
      const cls = 'er-ctx-item' + (item.danger ? ' danger' : '') + (item.disabled ? ' disabled' : '');
      const el = document.createElement('div');
      el.className = cls;
      el.innerHTML = (item.icon || '') + ' ' + escXml(item.label);
      el.onclick = () => { hideCtxMenu(); if (item.action) item.action(); };
      div.appendChild(el);
    });
    div.style.left = x + 'px'; div.style.top = y + 'px';
    document.body.appendChild(div);
    const r = div.getBoundingClientRect();
    if (r.right > window.innerWidth) div.style.left = (window.innerWidth - r.width - 4) + 'px';
    if (r.bottom > window.innerHeight) div.style.top = (window.innerHeight - r.height - 4) + 'px';
    ctxMenu = div;
  }

  document.addEventListener('mousedown', e => { if (ctxMenu && !ctxMenu.contains(e.target)) hideCtxMenu(); });
  document.addEventListener('keydown', e => { if (e.key === 'Escape') hideCtxMenu(); });

  // Always suppress browser context menu on ER canvas; our custom menu is triggered from mouseup
  canvas.addEventListener('contextmenu', e => {
    if (!erVisible || !erData) return;
    e.preventDefault();
  });

  function showErContextMenu(clientX, clientY) {
    const pt = screenToSvg(clientX, clientY);
    const isRo = typeof readOnly !== 'undefined' && readOnly;

    if (isRo) {
      showCtxMenu(clientX, clientY, [{ label: 'Read-only mode', disabled: true, icon: '&#128274;' }]);
      return;
    }

    const colHit = hitTestColumn(pt.x, pt.y);
    const tableHit = hitTestTable(pt.x, pt.y);

    if (colHit && tableHit) {
      const tbl = (erData.tables || []).find(t => t.name === colHit.table);
      const col = tbl && tbl.columns[colHit.idx];
      showCtxMenu(clientX, clientY, [
        { label: 'Edit Column...', icon: '&#9998;', action: () => showEditColumnDialog(colHit.table, col) },
        { label: 'Delete Column', icon: '&#128465;', danger: true, action: () => showDeleteColumnDialog(colHit.table, col) },
        '---',
        { label: 'Add Column...', icon: '&#10010;', action: () => showAddColumnDialog(colHit.table) },
        '---',
        { label: 'Rename Table...', icon: '&#9998;', action: () => showRenameTableDialog(colHit.table) },
        { label: 'Delete Table', icon: '&#128465;', danger: true, action: () => showDeleteTableDialog(colHit.table) },
      ]);
    } else if (tableHit) {
      showCtxMenu(clientX, clientY, [
        { label: 'Rename Table...', icon: '&#9998;', action: () => showRenameTableDialog(tableHit) },
        { label: 'Add Column...', icon: '&#10010;', action: () => showAddColumnDialog(tableHit) },
        '---',
        { label: 'Delete Table', icon: '&#128465;', danger: true, action: () => showDeleteTableDialog(tableHit) },
      ]);
    } else {
      showCtxMenu(clientX, clientY, [
        { label: 'Create Table...', icon: '&#10010;', action: () => showCreateTableDialog(pt) },
      ]);
    }
  }

  // ══════════════════════════════════════════════
  //  Table Structure Editor
  // ══════════════════════════════════════════════

  function showTableStructureDialog(tableName) {
    const db = erDB();
    const tbl = (erData.tables || []).find(t => t.name === tableName);
    if (!tbl) return;
    const cols = tbl.columns || [];

    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    function colRowHtml(col, i) {
      return '<tr data-orig="'+escXml(col.name)+'" data-idx="'+i+'">'
        + '<td style="padding:4px 6px">'
          + (col.pk ? '<span style="color:#e67700;font-size:11px" title="Primary Key">&#128273;</span>' : '')
        + '</td>'
        + '<td style="padding:4px 6px"><input type="text" value="'+escXml(col.name)+'" class="ts-name" style="width:100%;padding:4px 6px;border:1px solid #d0d3da;border-radius:3px;font-size:12px;font-family:monospace;outline:none"/></td>'
        + '<td style="padding:4px 6px"><select class="ts-type" style="width:100%;padding:4px 6px;border:1px solid #d0d3da;border-radius:3px;font-size:12px;font-family:monospace;outline:none">'
          + MYSQL_TYPES.map(t => '<option'+(t===col.type?' selected':'')+'>'+t+'</option>').join('')
        + '</select></td>'
        + '<td style="padding:4px 6px;text-align:center">'
          + '<button class="ts-del" title="Delete column" style="background:none;border:none;cursor:pointer;color:#ef4444;font-size:13px;padding:2px 4px">'
          + '<svg width="12" height="12" viewBox="0 0 24 24" fill="#ef4444"><path d="M6,7V19a2,2,0,0,0,2,2h8a2,2,0,0,0,2-2V7ZM8,3V2a1,1,0,0,1,1-1h6a1,1,0,0,1,1,1V3h4a1,1,0,0,1,0,2H4A1,1,0,0,1,4,3Z"/></svg>'
          + '</button>'
        + '</td>'
        + '</tr>';
    }

    let tableHtml = '<table style="width:100%;border-collapse:collapse;font-size:12px">'
      + '<thead><tr style="color:#6b7084;font-size:11px;text-align:left"><th style="padding:4px 6px;width:28px"></th><th style="padding:4px 6px">Name</th><th style="padding:4px 6px">Type</th><th style="padding:4px 6px;width:40px"></th></tr></thead>'
      + '<tbody id="tsCols">';
    cols.forEach((c, i) => { tableHtml += colRowHtml(c, i); });
    tableHtml += '</tbody></table>';

    overlay.innerHTML = '<div class="modal" style="width:560px;max-height:80vh">'
      + '<div class="modal-header"><span>Table: ' + escXml(tableName) + '</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body" style="overflow-y:auto">'
      + '<div class="form-field"><label>Table name</label>'
      + '<input type="text" id="tsTableName" value="'+escXml(tableName)+'" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/></div>'
      + '<div style="margin-bottom:8px;font-size:12px;font-weight:500;color:#6b7084">Columns</div>'
      + tableHtml
      + '<button class="btn btn-sm btn-secondary" id="tsAddCol" style="margin:8px 0 14px">+ Add Column</button>'
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="tsSqlPre" style="margin-top:8px;padding:10px;background:#f8f9fb;border-radius:4px;font-size:11px;border:1px solid #e0e2e8;white-space:pre-wrap;max-height:200px;overflow-y:auto"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="tsSaveBtn">Save Changes</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    let newColIdx = 0;

    document.getElementById('tsAddCol').onclick = () => {
      const tbody = document.getElementById('tsCols');
      const tr = document.createElement('tr');
      tr.setAttribute('data-new', 'true');
      tr.setAttribute('data-idx', 'new_' + (newColIdx++));
      tr.innerHTML = '<td style="padding:4px 6px"></td>'
        + '<td style="padding:4px 6px"><input type="text" placeholder="new_column" class="ts-name" style="width:100%;padding:4px 6px;border:1px solid #d0d3da;border-radius:3px;font-size:12px;font-family:monospace;outline:none"/></td>'
        + '<td style="padding:4px 6px"><select class="ts-type" style="width:100%;padding:4px 6px;border:1px solid #d0d3da;border-radius:3px;font-size:12px;font-family:monospace;outline:none">'
          + MYSQL_TYPES.map(t => '<option>'+t+'</option>').join('')
        + '</select></td>'
        + '<td style="padding:4px 6px;text-align:center">'
          + '<button class="ts-del" title="Remove" style="background:none;border:none;cursor:pointer;color:#ef4444;font-size:13px;padding:2px 4px">'
          + '<svg width="12" height="12" viewBox="0 0 24 24" fill="#ef4444"><path d="M6,7V19a2,2,0,0,0,2,2h8a2,2,0,0,0,2-2V7ZM8,3V2a1,1,0,0,1,1-1h6a1,1,0,0,1,1,1V3h4a1,1,0,0,1,0,2H4A1,1,0,0,1,4,3Z"/></svg>'
          + '</button>'
        + '</td>';
      tbody.appendChild(tr);
      tr.querySelector('.ts-name').focus();
      updateSql();
    };

    document.getElementById('tsCols').addEventListener('click', e => {
      const delBtn = e.target.closest('.ts-del');
      if (!delBtn) return;
      const tr = delBtn.closest('tr');
      if (tr.hasAttribute('data-new')) {
        tr.remove();
      } else {
        if (tr.classList.contains('ts-deleted')) {
          tr.classList.remove('ts-deleted');
          tr.style.opacity = '1'; tr.style.textDecoration = '';
        } else {
          tr.classList.add('ts-deleted');
          tr.style.opacity = '0.4'; tr.style.textDecoration = 'line-through';
        }
      }
      updateSql();
    });

    const pre = document.getElementById('tsSqlPre');

    function buildStmts() {
      const stmts = [];
      const newTableName = document.getElementById('tsTableName').value.trim();

      // MySQL: RENAME TABLE
      if (newTableName && newTableName !== tableName) {
        stmts.push('RENAME TABLE ' + quoteId(tableName) + ' TO ' + quoteId(newTableName));
      }
      const effectiveTable = (newTableName && newTableName !== tableName) ? newTableName : tableName;

      const rows = document.querySelectorAll('#tsCols tr');
      rows.forEach(tr => {
        const origName = tr.getAttribute('data-orig');
        const isNew = tr.hasAttribute('data-new');
        const isDeleted = tr.classList.contains('ts-deleted');
        const nameInput = tr.querySelector('.ts-name');
        const typeSelect = tr.querySelector('.ts-type');
        if (!nameInput || !typeSelect) return;
        const newName = nameInput.value.trim();
        const newType = typeSelect.value;

        if (isDeleted && origName) {
          stmts.push('ALTER TABLE ' + quoteId(effectiveTable) + ' DROP COLUMN ' + quoteId(origName));
        } else if (isNew && newName) {
          stmts.push('ALTER TABLE ' + quoteId(effectiveTable) + ' ADD COLUMN ' + quoteId(newName) + ' ' + newType);
        } else if (origName && !isNew && !isDeleted) {
          const origCol = cols.find(c => c.name === origName);
          if (origCol) {
            if (newName && newName !== origName) {
              stmts.push('ALTER TABLE ' + quoteId(effectiveTable) + ' RENAME COLUMN ' + quoteId(origName) + ' TO ' + quoteId(newName));
            }
            if (newType !== origCol.type) {
              const cn = (newName && newName !== origName) ? newName : origName;
              stmts.push('ALTER TABLE ' + quoteId(effectiveTable) + ' MODIFY COLUMN ' + quoteId(cn) + ' ' + newType);
            }
          }
        }
      });

      return stmts;
    }

    function updateSql() {
      const stmts = buildStmts();
      pre.textContent = stmts.length > 0 ? stmts.join(';\n') + ';' : '(no changes)';
    }
    overlay.addEventListener('input', updateSql);
    overlay.addEventListener('change', updateSql);
    updateSql();

    document.getElementById('tsSaveBtn').onclick = async function() {
      const stmts = buildStmts();
      if (stmts.length === 0) { overlay.remove(); return; }
      this.disabled = true; this.textContent = 'Saving...';
      try {
        await execStmts(stmts, db);
        const newTableName = document.getElementById('tsTableName').value.trim();
        if (newTableName && newTableName !== tableName && tablePositions[tableName]) {
          tablePositions[newTableName] = tablePositions[tableName];
          delete tablePositions[tableName];
          savePositions();
        }
        overlay.remove();
        showToast('Table updated', 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Save Changes';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Rename Table
  // ══════════════════════════════════════════════

  function showRenameTableDialog(tableName) {
    const db = erDB();
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    overlay.innerHTML = '<div class="modal" style="width:440px">'
      + '<div class="modal-header"><span>Rename Table</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<div class="form-field"><label>Current name</label>'
      + '<input type="text" value="'+escXml(tableName)+'" disabled style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;background:#f8f9fb;outline:none"/></div>'
      + '<div class="form-field"><label>New name</label>'
      + '<input type="text" id="erRenameInput" value="'+escXml(tableName)+'" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/></div>'
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erRenameSql" style="margin-top:8px;padding:10px;background:#f8f9fb;border-radius:4px;font-size:11px;border:1px solid #e0e2e8;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erRenameBtn">Rename</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    const input = document.getElementById('erRenameInput');
    const pre = document.getElementById('erRenameSql');
    function updateSql() {
      pre.textContent = 'RENAME TABLE ' + quoteId(tableName) + ' TO ' + quoteId(input.value.trim()) + ';';
    }
    input.addEventListener('input', updateSql); updateSql();
    input.focus(); input.select();

    document.getElementById('erRenameBtn').onclick = async function() {
      const newName = input.value.trim();
      if (!newName || newName === tableName) { overlay.remove(); return; }
      this.disabled = true; this.textContent = 'Renaming...';
      try {
        await execSql('RENAME TABLE ' + quoteId(tableName) + ' TO ' + quoteId(newName), db);
        if (tablePositions[tableName]) {
          tablePositions[newName] = tablePositions[tableName];
          delete tablePositions[tableName];
          savePositions();
        }
        overlay.remove();
        showToast('Table renamed: ' + tableName + ' \u2192 ' + newName, 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Rename';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Delete Table
  // ══════════════════════════════════════════════

  function showDeleteTableDialog(tableName) {
    const db = erDB();
    const rels = (erData.relations || []).filter(r => r.from_table === tableName || r.to_table === tableName);

    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    let warningHtml = '';
    if (rels.length > 0) {
      warningHtml = '<div style="padding:10px 14px;background:#fffbeb;border:1px solid #fde68a;border-radius:6px;margin-bottom:14px;font-size:12px;color:#92400e">'
        + '<strong>Warning:</strong> This table has ' + rels.length + ' foreign key relationship(s):<ul style="margin:6px 0 0;padding-left:18px">';
      rels.forEach(r => {
        warningHtml += '<li><code>' + escXml(r.from_table + '.' + r.from_column) + ' \u2192 ' + escXml(r.to_table + '.' + r.to_column) + '</code></li>';
      });
      warningHtml += '</ul></div>';
    }

    // MySQL doesn't support CASCADE on DROP TABLE the same way PG does (it's ignored),
    // but we keep the option for user awareness
    overlay.innerHTML = '<div class="modal" style="width:500px">'
      + '<div class="modal-header"><span>Delete Table</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<p style="margin:0 0 12px;font-size:13px">Drop table <strong style="font-family:monospace">' + escXml(tableName) + '</strong>?</p>'
      + warningHtml
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erDropSql" style="margin-top:8px;padding:10px;background:#fef2f2;border-radius:4px;font-size:11px;border:1px solid #fecaca;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erDropTableBtn" style="background:#ef4444;border-color:#ef4444">Drop Table</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    const pre = document.getElementById('erDropSql');
    pre.textContent = 'DROP TABLE ' + quoteId(tableName) + ';';

    document.getElementById('erDropTableBtn').onclick = async function() {
      this.disabled = true; this.textContent = 'Dropping...';
      try {
        await execSql('DROP TABLE ' + quoteId(tableName), db);
        delete tablePositions[tableName]; savePositions();
        overlay.remove();
        showToast('Table dropped: ' + tableName, 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Drop Table';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Create Table
  // ══════════════════════════════════════════════

  function showCreateTableDialog(pos) {
    const db = erDB();
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    const colRowHtml = function(i, name, type, pk) {
      return '<div class="er-ct-row" data-idx="'+i+'" style="display:flex;gap:6px;align-items:center;margin-bottom:6px">'
        + '<input type="text" placeholder="column_name" value="'+(name||'')+'" style="flex:1;padding:5px 8px;border:1px solid #d0d3da;border-radius:4px;font-size:12px;font-family:monospace;outline:none" class="er-ct-name">'
        + '<select style="width:130px;padding:5px 6px;border:1px solid #d0d3da;border-radius:4px;font-size:12px;font-family:monospace;outline:none" class="er-ct-type">'
        + MYSQL_TYPES.map(t => '<option'+(t===type?' selected':'')+'>'+t+'</option>').join('')
        + '</select>'
        + '<label style="font-size:11px;white-space:nowrap;cursor:pointer"><input type="checkbox" class="er-ct-pk" '+(pk?'checked':'')+'>PK</label>'
        + '<button onclick="this.parentElement.remove()" style="background:none;border:none;cursor:pointer;color:#ef4444;font-size:14px" title="Remove">&times;</button>'
        + '</div>';
    };

    overlay.innerHTML = '<div class="modal" style="width:560px">'
      + '<div class="modal-header"><span>Create Table</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<div class="form-field"><label>Table name</label>'
      + '<input type="text" id="erCtName" placeholder="new_table" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/></div>'
      + '<div style="margin-bottom:8px;font-size:12px;font-weight:500;color:#6b7084">Columns</div>'
      + '<div id="erCtCols">' + colRowHtml(0, 'id', 'int', true) + '</div>'
      + '<button class="btn btn-sm btn-secondary" id="erCtAddCol" style="margin-bottom:14px">+ Add Column</button>'
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erCtSql" style="margin-top:8px;padding:10px;background:#f8f9fb;border-radius:4px;font-size:11px;border:1px solid #e0e2e8;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erCtBtn">Create</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    let colIdx = 1;
    document.getElementById('erCtAddCol').onclick = () => {
      document.getElementById('erCtCols').insertAdjacentHTML('beforeend', colRowHtml(colIdx++, '', 'varchar(255)', false));
      updateSql();
    };

    const pre = document.getElementById('erCtSql');
    function buildSql() {
      const name = document.getElementById('erCtName').value.trim();
      if (!name) return '';
      const rows = document.querySelectorAll('.er-ct-row');
      const cols = []; const pks = [];
      rows.forEach(r => {
        const cn = r.querySelector('.er-ct-name').value.trim();
        const ct = r.querySelector('.er-ct-type').value;
        const cp = r.querySelector('.er-ct-pk').checked;
        if (cn) {
          let colDef = '  ' + quoteId(cn) + ' ' + ct;
          // Auto-add AUTO_INCREMENT for int PK
          if (cp && /int/i.test(ct)) colDef += ' AUTO_INCREMENT';
          cols.push(colDef);
          if (cp) pks.push(quoteId(cn));
        }
      });
      if (cols.length === 0) return '';
      let sql = 'CREATE TABLE ' + quoteId(name) + ' (\n' + cols.join(',\n');
      if (pks.length > 0) sql += ',\n  PRIMARY KEY (' + pks.join(', ') + ')';
      sql += '\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;';
      return sql;
    }
    function updateSql() { pre.textContent = buildSql(); }
    overlay.addEventListener('input', updateSql);
    overlay.addEventListener('change', updateSql);
    updateSql();
    document.getElementById('erCtName').focus();

    document.getElementById('erCtBtn').onclick = async function() {
      const sql = buildSql();
      if (!sql) { showToast('Table name and at least one column required', 'warning'); return; }
      this.disabled = true; this.textContent = 'Creating...';
      try {
        await execSql(sql, db);
        const tblName = document.getElementById('erCtName').value.trim();
        // Pre-inject position into localStorage so loadER restores it at the right-click point
        if (pos) {
          try {
            const raw = localStorage.getItem(storageKey());
            const saved = raw ? JSON.parse(raw) : { v:1, layout: currentLayout, tables:{}, view:{x:viewX,y:viewY,scale:scale} };
            saved.tables[tblName] = { x: Math.round(pos.x), y: Math.round(pos.y) };
            localStorage.setItem(storageKey(), JSON.stringify(saved));
          } catch(e) {}
        }
        overlay.remove();
        showToast('Table created: ' + tblName, 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Create';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Add Column
  // ══════════════════════════════════════════════

  function showAddColumnDialog(tableName) {
    const db = erDB();
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    overlay.innerHTML = '<div class="modal" style="width:480px">'
      + '<div class="modal-header"><span>Add Column to ' + escXml(tableName) + '</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<div class="form-field"><label>Column name</label>'
      + '<input type="text" id="erAcName" placeholder="column_name" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/></div>'
      + '<div class="form-field"><label>Type</label>' + mysqlTypeSelect('varchar(255)', 'erAcType') + '</div>'
      + '<div style="display:flex;gap:16px;margin-bottom:14px">'
      + '<label style="font-size:12px;cursor:pointer"><input type="checkbox" id="erAcNotNull"> NOT NULL</label>'
      + '</div>'
      + '<div class="form-field"><label>Default value (optional)</label>'
      + '<input type="text" id="erAcDefault" placeholder="" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/></div>'
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erAcSql" style="margin-top:8px;padding:10px;background:#f8f9fb;border-radius:4px;font-size:11px;border:1px solid #e0e2e8;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erAcBtn">Add</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    const pre = document.getElementById('erAcSql');
    function buildSql() {
      const name = document.getElementById('erAcName').value.trim();
      if (!name) return '';
      let sql = 'ALTER TABLE ' + quoteId(tableName) + ' ADD COLUMN ' + quoteId(name) + ' ' + document.getElementById('erAcType').value;
      if (document.getElementById('erAcNotNull').checked) sql += ' NOT NULL';
      const def = document.getElementById('erAcDefault').value.trim();
      if (def) sql += ' DEFAULT ' + def;
      return sql;
    }
    function updateSql() { pre.textContent = buildSql() ? buildSql() + ';' : ''; }
    overlay.addEventListener('input', updateSql);
    overlay.addEventListener('change', updateSql);
    updateSql();
    document.getElementById('erAcName').focus();

    document.getElementById('erAcBtn').onclick = async function() {
      const sql = buildSql();
      if (!sql) { showToast('Column name is required', 'warning'); return; }
      this.disabled = true; this.textContent = 'Adding...';
      try {
        await execSql(sql, db);
        overlay.remove();
        showToast('Column added to ' + tableName, 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Add';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Edit Column
  // ══════════════════════════════════════════════

  function showEditColumnDialog(tableName, col) {
    if (!col) return;
    const db = erDB();
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    const colType = col.type || 'varchar(255)';

    overlay.innerHTML = '<div class="modal" style="width:480px">'
      + '<div class="modal-header"><span>Edit Column: ' + escXml(tableName) + '.' + escXml(col.name) + '</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<div class="form-field"><label>Column name</label>'
      + '<input type="text" id="erEcName" value="'+escXml(col.name)+'" style="width:100%;padding:7px 10px;border:1px solid #d0d3da;border-radius:5px;font-size:13px;font-family:monospace;outline:none"/></div>'
      + '<div class="form-field"><label>Type</label>' + mysqlTypeSelect(colType, 'erEcType') + '</div>'
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erEcSql" style="margin-top:8px;padding:10px;background:#f8f9fb;border-radius:4px;font-size:11px;border:1px solid #e0e2e8;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erEcBtn">Save</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    const pre = document.getElementById('erEcSql');
    function buildStmts() {
      const newName = document.getElementById('erEcName').value.trim();
      const newType = document.getElementById('erEcType').value;
      const stmts = [];
      if (newName && newName !== col.name) {
        stmts.push('ALTER TABLE ' + quoteId(tableName) + ' RENAME COLUMN ' + quoteId(col.name) + ' TO ' + quoteId(newName));
      }
      if (newType && newType !== colType) {
        const cn = newName && newName !== col.name ? newName : col.name;
        stmts.push('ALTER TABLE ' + quoteId(tableName) + ' MODIFY COLUMN ' + quoteId(cn) + ' ' + newType);
      }
      return stmts;
    }
    function updateSql() {
      const stmts = buildStmts();
      pre.textContent = stmts.length > 0 ? stmts.join(';\n') + ';' : '(no changes)';
    }
    overlay.addEventListener('input', updateSql);
    overlay.addEventListener('change', updateSql);
    updateSql();
    document.getElementById('erEcName').focus();

    document.getElementById('erEcBtn').onclick = async function() {
      const stmts = buildStmts();
      if (stmts.length === 0) { overlay.remove(); return; }
      this.disabled = true; this.textContent = 'Saving...';
      try {
        await execStmts(stmts, db);
        overlay.remove();
        showToast('Column updated', 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Save';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Delete Column
  // ══════════════════════════════════════════════

  function showDeleteColumnDialog(tableName, col) {
    if (!col) return;
    const db = erDB();
    const rels = (erData.relations || []).filter(r =>
      (r.from_table === tableName && r.from_column === col.name) ||
      (r.to_table === tableName && r.to_column === col.name)
    );

    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay'; overlay.style.zIndex = '400';

    let warningHtml = '';
    if (col.pk) {
      warningHtml += '<div style="padding:10px 14px;background:#fef2f2;border:1px solid #fecaca;border-radius:6px;margin-bottom:10px;font-size:12px;color:#991b1b">'
        + '<strong>Warning:</strong> This is a PRIMARY KEY column.</div>';
    }
    if (rels.length > 0) {
      warningHtml += '<div style="padding:10px 14px;background:#fffbeb;border:1px solid #fde68a;border-radius:6px;margin-bottom:10px;font-size:12px;color:#92400e">'
        + '<strong>Warning:</strong> This column has ' + rels.length + ' foreign key relationship(s):<ul style="margin:6px 0 0;padding-left:18px">';
      rels.forEach(r => { warningHtml += '<li><code>' + escXml(r.from_table+'.'+r.from_column+' \u2192 '+r.to_table+'.'+r.to_column) + '</code></li>'; });
      warningHtml += '</ul></div>';
    }

    // MySQL doesn't support CASCADE on DROP COLUMN
    overlay.innerHTML = '<div class="modal" style="width:480px">'
      + '<div class="modal-header"><span>Delete Column</span>'
      + '<button class="modal-close" onclick="this.closest(\'.modal-overlay\').remove()">&times;</button></div>'
      + '<div class="modal-body">'
      + '<p style="margin:0 0 12px;font-size:13px">Drop column <strong style="font-family:monospace">' + escXml(tableName) + '.' + escXml(col.name) + '</strong> <span style="color:#94a3b8">(' + escXml(col.type) + ')</span>?</p>'
      + warningHtml
      + '<details><summary style="cursor:pointer;font-size:12px;color:#6b7084">Preview SQL</summary>'
      + '<pre id="erDcSql" style="margin-top:8px;padding:10px;background:#fef2f2;border-radius:4px;font-size:11px;border:1px solid #fecaca;white-space:pre-wrap"></pre></details>'
      + '</div>'
      + '<div class="modal-footer">'
      + '<button class="btn btn-secondary" onclick="this.closest(\'.modal-overlay\').remove()">Cancel</button>'
      + '<button class="btn btn-primary" id="erDcBtn" style="background:#ef4444;border-color:#ef4444">Drop Column</button>'
      + '</div></div>';
    document.body.appendChild(overlay);

    const pre = document.getElementById('erDcSql');
    pre.textContent = 'ALTER TABLE ' + quoteId(tableName) + ' DROP COLUMN ' + quoteId(col.name) + ';';

    document.getElementById('erDcBtn').onclick = async function() {
      this.disabled = true; this.textContent = 'Dropping...';
      try {
        await execSql('ALTER TABLE ' + quoteId(tableName) + ' DROP COLUMN ' + quoteId(col.name), db);
        overlay.remove();
        showToast('Column dropped: ' + tableName + '.' + col.name, 'success');
        collabNotifySchemaChanged();
        loadER();
      } catch(e) {
        showToast('Error: ' + e.message, 'error');
        this.disabled = false; this.textContent = 'Drop Column';
      }
    };
  }

  // ══════════════════════════════════════════════
  //  Collaboration — WebSocket client
  // ══════════════════════════════════════════════

  function collabConnect() {
    if (collabWs && collabWs.readyState <= 1) return;
    const db = erCurrentDB();
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const url = proto + '//' + location.host + '/ws/er?db=' + encodeURIComponent(db);
    try {
      collabWs = new WebSocket(url);
    } catch(e) { return; }

    collabWs.onopen = function() {
      clearTimeout(collabReconnectTimer);
    };

    collabWs.onmessage = function(evt) {
      try {
        const msg = JSON.parse(evt.data);
        collabHandleMsg(msg);
      } catch(e) {}
    };

    collabWs.onclose = function() {
      collabWs = null;
      collabCursors = {};
      collabDrags = {};
      collabUsers = [];
      collabRenderPresence();
      if (erVisible) {
        collabReconnectTimer = setTimeout(collabConnect, 2000);
      }
    };

    collabWs.onerror = function() {};
  }

  function collabDisconnect() {
    clearTimeout(collabReconnectTimer);
    if (collabWs) {
      collabWs.onclose = null;
      collabWs.close();
      collabWs = null;
    }
    collabMe = null;
    collabUsers = [];
    collabCursors = {};
    collabDrags = {};
    collabRenderPresence();
  }

  function collabSend(msg) {
    if (collabWs && collabWs.readyState === 1) {
      collabWs.send(JSON.stringify(msg));
    }
  }

  function collabHandleMsg(msg) {
    switch (msg.type) {
      case 'welcome':
        collabMe = { user_id: msg.user_id, name: msg.name, color: msg.color };
        break;

      case 'presence':
        if (msg.data) {
          collabUsers = msg.data;
          collabRenderPresence();
        }
        break;

      case 'cursor':
        if (msg.user_id && msg.user_id !== (collabMe && collabMe.user_id)) {
          const d = msg.data ? JSON.parse(msg.data) : null;
          if (d) {
            collabCursors[msg.user_id] = {
              x: d.x, y: d.y, name: msg.name, color: msg.color, ts: Date.now()
            };
            scheduleRender();
          }
        }
        break;

      case 'drag':
        if (msg.user_id && msg.user_id !== (collabMe && collabMe.user_id)) {
          const d = msg.data ? JSON.parse(msg.data) : null;
          if (d) {
            if (d.done) {
              delete collabDrags[msg.user_id];
              if (d.table && tablePositions[d.table]) {
                tablePositions[d.table] = { x: d.x, y: d.y };
                savePositions();
              }
            } else {
              collabDrags[msg.user_id] = {
                table: d.table, x: d.x, y: d.y, name: msg.name, color: msg.color
              };
              if (d.table && tablePositions[d.table]) {
                tablePositions[d.table] = { x: d.x, y: d.y };
              }
            }
            scheduleRender();
          }
        }
        break;

      case 'schema_changed':
        loadER();
        break;
    }
  }

  let lastCursorSend = 0;
  function collabSendCursor(svgX, svgY) {
    const now = Date.now();
    if (now - lastCursorSend < 50) return;
    lastCursorSend = now;
    collabSend({ type: 'cursor', data: JSON.stringify({ x: svgX, y: svgY }) });
  }

  let lastDragSend = 0;
  function collabSendDrag(tableName, x, y, done) {
    const now = Date.now();
    if (!done && now - lastDragSend < 33) return;
    lastDragSend = now;
    collabSend({ type: 'drag', data: JSON.stringify({ table: tableName, x: x, y: y, done: !!done }) });
  }

  function collabNotifySchemaChanged() {
    collabSend({ type: 'schema_changed' });
  }

  window._collabSwitchDB = collabSwitchDB;
  function collabSwitchDB(newDB) {
    collabSend({ type: 'switch_db', data: JSON.stringify({ db: newDB }) });
    collabDisconnect();
    setTimeout(collabConnect, 100);
  }

  function collabRenderPresence() {
    let bar = document.getElementById('erCollabBar');
    if (!bar) return;

    const others = collabUsers.filter(u => !collabMe || u.user_id !== collabMe.user_id);
    if (others.length === 0 && !collabMe) {
      bar.innerHTML = '';
      return;
    }

    let html = '';
    if (collabMe) {
      html += '<div class="collab-avatar" style="background:' + collabMe.color + '" title="You (' + escXml(collabMe.name) + ')">'
        + escXml(collabMe.name.charAt(0)) + '</div>';
    }
    others.forEach(u => {
      html += '<div class="collab-avatar" style="background:' + u.color + '" title="' + escXml(u.name) + '">'
        + escXml(u.name.charAt(0)) + '</div>';
    });
    const total = collabUsers.length;
    if (total > 0) {
      html += '<span class="collab-count">' + total + ' online</span>';
    }
    bar.innerHTML = html;
  }

  function renderCollabOverlays() {
    let s = '';
    const now = Date.now();

    for (const uid in collabCursors) {
      const c = collabCursors[uid];
      if (now - c.ts > 5000) { delete collabCursors[uid]; continue; }
      const opacity = now - c.ts > 3000 ? '0.4' : '0.85';
      s += '<g opacity="' + opacity + '">'
        + '<polygon points="' + c.x + ',' + c.y + ' ' + (c.x) + ',' + (c.y+16) + ' ' + (c.x+5) + ',' + (c.y+13) + ' ' + (c.x+10) + ',' + (c.y+19) + ' ' + (c.x+13) + ',' + (c.y+17) + ' ' + (c.x+8) + ',' + (c.y+11) + ' ' + (c.x+14) + ',' + (c.y+10) + '" '
        + 'fill="' + c.color + '" stroke="white" stroke-width="0.8"/>'
        + '<rect x="' + (c.x+14) + '" y="' + (c.y+10) + '" width="' + (c.name.length * 7 + 10) + '" height="18" rx="9" fill="' + c.color + '"/>'
        + '<text x="' + (c.x+19) + '" y="' + (c.y+23) + '" font-size="11" fill="white" font-family="-apple-system,BlinkMacSystemFont,sans-serif" font-weight="500">'
        + escXml(c.name) + '</text>'
        + '</g>';
    }

    for (const uid in collabDrags) {
      const d = collabDrags[uid];
      const pos = tablePositions[d.table];
      if (!pos) continue;
      const tbl = erData && (erData.tables || []).find(t => t.name === d.table);
      if (!tbl) continue;
      const h = getTableHeight(tbl);
      const w = tw(d.table);
      s += '<rect x="' + (pos.x-3) + '" y="' + (pos.y-3) + '" width="' + (w+6) + '" height="' + (h+6)
        + '" rx="8" fill="none" stroke="' + d.color + '" stroke-width="2.5" stroke-dasharray="6,3" opacity="0.7"/>'
        + '<rect x="' + (pos.x-3) + '" y="' + (pos.y-22) + '" width="' + (d.name.length * 7 + 10) + '" height="18" rx="9" fill="' + d.color + '"/>'
        + '<text x="' + (pos.x+2) + '" y="' + (pos.y-9) + '" font-size="11" fill="white" font-family="-apple-system,BlinkMacSystemFont,sans-serif" font-weight="500">'
        + escXml(d.name) + '</text>';
    }

    return s;
  }

  function collabStartCleanup() {
    clearInterval(collabCursorCleanTimer);
    collabCursorCleanTimer = setInterval(() => {
      const now = Date.now();
      let changed = false;
      for (const uid in collabCursors) {
        if (now - collabCursors[uid].ts > 5000) { delete collabCursors[uid]; changed = true; }
      }
      if (changed) scheduleRender();
    }, 2000);
  }

})();
