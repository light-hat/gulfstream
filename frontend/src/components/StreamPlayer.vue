<template>
  <div>
    <video ref="videoPlayer" controls autoplay></video>
    <div>
      <select v-model="selectedQuality">
        <option v-for="quality in qualities" :key="quality" :value="quality">
          {{ quality }}
        </option>
      </select>
    </div>
  </div>
</template>

<script>
import Hls from 'hls.js';

export default {
  props: ['streamId'],
  data() {
    return {
      hls: null,
      qualities: [],
      selectedQuality: null
    }
  },
  mounted() {
    this.initPlayer();
  },
  watch: {
    selectedQuality(newQuality) {
      this.changeQuality(newQuality);
    }
  },
  methods: {
    async initPlayer() {
      const masterUrl = `http://hls-origin:8080/${this.streamId}/master.m3u8`;
      
      // Получение доступных качеств
      const response = await fetch(masterUrl);
      const playlist = await response.text();
      this.parseQualities(playlist);
      
      // Инициализация плеера
      if (Hls.isSupported()) {
        this.hls = new Hls({
          enableWorker: true,
          lowLatencyMode: true
        });
        
        this.hls.loadSource(masterUrl);
        this.hls.attachMedia(this.$refs.videoPlayer);
        
        this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
          this.$refs.videoPlayer.play();
        });
      } else if (this.$refs.videoPlayer.canPlayType('application/vnd.apple.mpegurl')) {
        this.$refs.videoPlayer.src = masterUrl;
      }
    },
    parseQualities(playlist) {
      const regex = /STREAM\-INF.*?RESOLUTION=(\d+x\d+).*?\n(.*?\.m3u8)/g;
      let match;
      
      while ((match = regex.exec(playlist)) !== null) {
        this.qualities.push({
          resolution: match[1],
          url: match[2]
        });
      }
      
      if (this.qualities.length > 0) {
        this.selectedQuality = this.qualities[0].url;
      }
    },
    changeQuality(qualityUrl) {
      if (this.hls) {
        this.hls.loadSource(qualityUrl);
      }
    }
  },
  beforeDestroy() {
    if (this.hls) {
      this.hls.destroy();
    }
  }
}
</script>
