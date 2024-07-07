var indexData = null

function loadIndex() {
  return fetch("/index.json")
    .then(response => response.json())
    .then(data => {
      indexData = data
    })
}

// tagged("tag") returns a list of documents with that tag.
function tagged(tag) {
  const t = indexData.tm[tag]
  if (t === undefined) {
    return []
  }
  var hits = []
  for (const id of t) {
    hits.push(indexData.docs[id])
  }
  return hits
}

// search("query") returns a list of documents that match the query.
function search(query) {

  query = query.toLowerCase()

  // split query by space
  var parts = query.split(" ")

  var hitMap = {}
  for (var p of parts) {
    const pmap = indexData.kw[p]
    if (pmap === undefined) {
      continue // no hits
    }
    for (const pid of Object.keys(pmap)) {
      if (hitMap[pid] === undefined) {
        hitMap[pid] = 0
      }
      hitMap[pid] += pmap[pid]
    }
  }

  var sortedHitIds = Object.keys(hitMap).sort(function (a, b) {
    return hitMap[b] - hitMap[a]
  })

  var hits = []
  for (var i = 0; i < sortedHitIds.length; i++) {
    var id = sortedHitIds[i]
    var hit = indexData.docs[id]
    hits.push(hit)
  }
  return hits
}