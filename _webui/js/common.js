







function videoStream(vi) {
  for (var i=0; i<vi.ffprobe.streams.length; i++) {
    var stream = vi.ffprobe.streams[i];
    if (stream.codec_type="video") {
      return stream;
    }
  }
}

// this is probably useless, need audioStreams()
function audioStream(vi) {
  for (var i=0; i<vi.ffprobe.streams.length; i++) {
    var stream = vi.ffprobe.streams[i];
    if (stream.codec_type="audio") {
      return stream;
    }
  }
}