<template>
  <div style="height:100px;">
    <button
      class="button is-info"
      v-show="isShowNewVisible"
      v-on:click="getNewLogItems()"
    >
      Show new
    </button>
    <div style="height:10px; overflow=scrollable;">
      <table>
        <tr v-for="(item, index) in items" v-bind:key="index">
          <td style="width:165px">
            <p>[{{ item.atd }}&nbsp;{{ item.att }}]</p>
          </td>
          <td
            :style="{
              'background-color': config.severityColor[item.severity],
            }"
          >
            [{{ config.severityLabel[item.severity] }}]
          </td>
          <td>{{ item.msg }}</td>
        </tr>
      </table>

      <button
        class="button"
        v-show="isShowMoreVisible"
        v-on:click="getNextLogItems(config.logPerPage)"
      >
        Show previous
      </button>
      <button
        class="button"
        v-show="isShowMoreVisible"
        v-on:click="getAllLogItems()"
      >
        Show All
      </button>
    </div>
  </div>
</template>

<script>
function formatD(date) {
  var d = new Date(date),
    month = "" + (d.getMonth() + 1),
    day = "" + d.getDate(),
    year = d.getFullYear();

  if (month < 10) month = "0" + month;
  if (day < 10) day = "0" + day;

  return [year, month, day].join("-");
}

function formatT(date) {
  var d = new Date(date),
    h = d.getHours(),
    m = d.getMinutes(),
    s = d.getSeconds(),
    ms = d.getMilliseconds();

  if (h < 10) h = "0" + h;
  if (m < 10) m = "0" + m;
  if (s < 10) s = "0" + s;
  if (ms < 10) {
    ms = "00" + ms;
  }
  if (ms >= 10 && ms < 100) {
    ms = "0" + ms;
  }

  return [h, m, s].join(":") + "." + ms;
}

async function fetchLogItems(id, from, to, filter) {
  let url = "/api/log/" + id + "/feed?from=" + from + "&to=" + to;
  if (filter) {
    url += "&filter=" + filter;
  }

  const resp = await fetch(url);
  var data = await resp.json();

  let res = [];
  for (var i in data) {
    //    console.log("tmp.i=", i);
    // console.log("r: ", res[i]);
    res.push({
      atd: formatD(data[i].at),
      att: formatT(data[i].at),
      msg: data[i].msg,
      severity: data[i].severity,
    });
    //res.at = formatDate(res[i].at);
    //console.log("res.lentht in for: ", res.length);
  }
  // console.log("res: ", res);

  return res;
}

export default {
  name: "InstanceModalLog",
  props: {
    code: String,
    filter: {
      type: String,
      default: "",
    },
  },
  data: function() {
    return {
      id: "",
      items: [],
      isShowMoreVisible: false,
      isShowNewVisible: false,
      totalLen: 0,
      newTotalLen: 0,
      polling: null,
      config: {
        logPerPage: 20,
        severityColor: {
          0: "red",
          1: "orange",
          2: "yellow",
          3: "",
        },
        severityLabel: {
          0: "FATL",
          1: "ERRO",
          2: "WARN",
          3: "INFO",
        },
      },
    };
  },
  computed: {
    sevcolor: function() {
      return "red";
    },
  },
  methods: {
    getLogItems: function(id) {
      let that = this;

      this.id = id;
      this.isShowMoreVisible = false;

      this.getCount(id).then((cnt) => {
        that.totalLen = cnt;

        fetchLogItems(id, 0, this.config.logPerPage - 1, this.filter).then(
          (arr) => {
            that.items.push(...arr);
            that.isShowMoreVisible = that.items.length < that.totalLen;
          }
        );

        that.polling = setInterval(function() {
          that.getCount(id).then((cnt) => {
            that.newTotalLen = cnt;
            if (that.newTotalLen > that.totalLen) {
              that.isShowNewVisible = true;
            }
          });

          /*console.log(
            "newTotalLen",
            that.newTotalLen,
            "totalLen",
            that.totalLen
          );*/
        }, 5000);
      });

      //console.log("Entering getLogItems(id), totalLen=", id, that.totalLen);
    },

    getNextLogItems: function(cnt) {
      let fincnt = -1;
      if (cnt != -1) {
        fincnt = this.items.length + cnt - 1;
      }
      fetchLogItems(this.id, this.items.length, fincnt).then((v) =>
        this.items.push(...v)
      );
    },

    getNewLogItems: function() {
      fetchLogItems(this.id, 0, this.newTotalLen - this.totalLen).then(
        (arr) => {
          this.items.unshift(...arr);
          this.totalLen = this.newTotalLen;
          this.isShowNewVisible = false;
        }
      );
    },

    getAllLogItems: function() {
      let that = this;
      fetchLogItems(this.id, 0, -1).then((arr) => {
        that.items = [];
        that.items.push(...arr);
        that.totalLen = arr.length;
        that.isShowMoreVisible = false;
      });
    },

    // getCount return total length of the list L:XXXX where XXXX is instance code.
    getCount: function(id) {
      return fetch("/api/log/" + id + "/count")
        .then((r) => r.json())
        .then((j) => {
          return j.count;
        });
    },
  },
  mounted: function() {
    //console.log("before updated instance modal info");
    this.getLogItems(this.$props.code);
  },
  beforeDestroy() {
    clearInterval(this.polling);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
tr:nth-child(even) {
  background-color: #f2f2f2;
}
</style>
