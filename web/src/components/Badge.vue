<template>
  <div class="badge">
    <div
      class="topRight"
      :style="{
        'top': -(size/2) + 'px',
        'right': -(size/2) + 'px',
        'height': size + 'px',
        'width': size + 'px',
      }"
    >
      <StatusIcon
        circle
        spinner
        class="readyWrap"
        :type="type"
        :size="size-4"
      >
      </StatusIcon>
    </div>

  </div>

</template>

<script>
import colors from "vuetify/lib/util/colors";
import StatusIcon from "./StatusIcon.vue";

export default {
  name: "Badge",
  components: {
    StatusIcon,
  },
  props: {
    type: String,
    blinking: Boolean,
    size: {
      default: 16,
      type: Number,
    },
  },
  methods: {
    getColor() {
      switch (this.type) {
        case "ready":
          return colors.green.base;
        case "error":
          return colors.red.accent2;
        case "loading":
          return colors.orange.darken1;
        default:
          return colors.grey;
      }
    },
  },
};
</script>

<style lang="less" scoped>
.blink-transition {
  transition: all 1s ease;
  opacity: 0%;
}

.blink-enter,
.blink-leave {
  opacity: 100%;
}

.badge .topRight {
  position: absolute;
  display: inline-block;
}

.badge .readyWrap {
  position: absolute;
  border: 2px solid #f8f3f2 !important;
  // box-shadow: 0px 0px 10px rgba(0, 0, 0, 1);
}

.readyIcon {
  display: inline-block;
  vertical-align: middle;
  text-align: center;
  // border-radius: 50%;
  // padding: 2px;
}

@-webkit-keyframes Blinking {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.2;
  }
  100% {
    opacity: 1;
  }
}
@-moz-keyframes Blinking {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.2;
  }
  100% {
    opacity: 1;
  }
}
@keyframes Blinking {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.2;
  }
  100% {
    opacity: 1;
  }
}

.blinking {
  -webkit-animation: Blinking 3s ease-in-out infinite !important;
  -moz-animation: Blinking 3s ease-in-out infinite !important;
  animation: Blinking 3s ease-in-out infinite !important;
}

// .readySpinner {
//   display: inline-block !important;
//   margin-bottom: 1px;
// }
</style>