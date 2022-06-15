<template>
  <div>
    <div
      :class="['readyWrap', {blinking: blinking}]"
      :style="{
        'background-color': getColor(),
        'top': -(size/2) + 'px',
        'right': -(size/2) + 'px',
        'width': (size-4) + 'px',
        'height': (size-4) + 'px',
      }"
    >

      <v-icon
        v-if="type==='ready'"
        class="readyIcon"
        color="white"
        :size="size-6"
      > mdi-check</v-icon>
      <v-icon
        v-else-if="type==='error'"
        class="readyIcon"
        color="white"
        :size="size-6"
      > mdi-exclamation</v-icon>
      <v-progress-circular
        v-else-if="type==='loading'"
        class="readySpinner"
        indeterminate
        :size="size-8"
        :width="1.5"
        color="white"
      >
      </v-progress-circular>
    </div>

  </div>

</template>

<script>
import colors from "vuetify/lib/util/colors";

export default {
  name: "Badge",
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
.readyButton {
  width: 20px;
  height: 20px;
  line-height: 18px;
}

.readyWrap {
  position: absolute;
  display: flex; // make us of Flexbox
  align-items: center; // does vertically center the desired content
  justify-content: center; // horizontally centers single line items
  text-align: center; // optional, but helps horizontally center text that breaks into multiple lines

  border-radius: 50%;
  border: 2px solid #f8f3f2;
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