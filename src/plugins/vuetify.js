import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';

import colors from 'vuetify/lib/util/colors'

Vue.use(Vuetify);

export default new Vuetify({
theme: {
  themes: {
    light: {
      primary: colors.blue.darken1, 
      secondary: colors.blue.lighten3, 
      accent: colors.lightBlue.base,
    },
  },
},
});
