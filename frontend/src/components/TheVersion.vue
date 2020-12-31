<template>
  <div>
    <div
      style="margin-right:3px; margin-left:1px; border-top:1px solid gray; margin-top:10px; padding-top:6px"
    >
      <span class="version-summary"
        >Peregrine Traders, 2020 ©<br />
        <span v-if="platformBuild != ''">Platform {{ platformBuild }},</span>
        GUI
        <a
          @click="
            {
              getVersion();
              isActive = true;
            }
          "
          >{{ version }}</a
        >
        built at {{ buildAt }} for
        {{ buildFor }}
      </span>
    </div>
    <responsive-modal
      name="window-version"
      :isActive.sync="isActive"
      title="Release Notes"
      widthLimit="80%"
    >
      <markdown-it-vue
        class="md-body"
        style="padding:0px 0px 10px 10px;overflow-y:auto; display:block; height: 93%; scroll:hidden; color:silver; text-align: left"
        :content="content"
        :options="options"
      />
    </responsive-modal>
  </div>
</template>

<script>
import MarkdownItVue from "markdown-it-vue";
import "markdown-it-vue/dist/markdown-it-vue.css";
import ResponsiveModal from "@/components/ResponsiveModal.vue";

export default {
  name: "Version",
  components: {
    MarkdownItVue,
    "responsive-modal": ResponsiveModal,
  },
  data: function() {
    return {
      version: "",
      buildAt: "",
      buildFor: "",
      content: "",
      options: { typographer: true },
      isActive: false,
    };
  },
  computed: {
    platformBuild: function() {
      return this.$store.getters.platformBuild;
    },
  },
  methods: {
    getVersion: function(withoutReleaseNotes) {
      let that = this;
      let url = "/api/version";
      if (withoutReleaseNotes) {
        url += "?withoutReleaseNotes=true";
      }

      fetch(url)
        .then((r) => r.json())
        .then(function(j) {
          that.isUpdated = that.version != "" && that.version != j.version;
          that.version = j.version;
          that.buildAt = j.buildAt;
          that.buildFor = j.buildFor;
          if (
            j.releaseNotes != undefined &&
            j.releaseNotes != null &&
            j.releaseNotes != ""
          ) {
            that.content = j.releaseNotes;
          }
        });
    },
  },
  mounted: function() {
    this.getVersion(true);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.version-summary {
  color: silver;
  vertical-align: middle;
  font-size: 8pt;
}

.version-summary:hover {
  color: white;
}
</style>
