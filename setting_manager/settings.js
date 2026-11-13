function initSettingManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.settingManagerUrls || {},
        settings: [],
        loading: false,
        error: null,
        savingKey: null,
        currencies: [
          'USD', 'EUR', 'GBP', 'JPY', 'AUD', 'CAD', 'CHF', 'CNY',
          'SEK', 'NZD', 'INR', 'BRL', 'ZAR', 'PLN', 'BGN', 'RON'
        ]
      };
    },
    mounted() {
      this.loadSettings();
    },
    methods: {
      async loadSettings() {
        this.loading = true;
        this.error = null;
        try {
          const response = await fetch(this.urls.loadSettings, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              page: 0,
              per_page: 1000,
              conditions: [],
              sort_by: 'setting_key',
              sort: 'asc'
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            if (result.data.urls) {
              this.urls = { ...this.urls, ...result.data.urls };
            }
            this.settings = result.data.settings || [];
          } else {
            this.error = result.message || 'Failed to load settings';
          }
        } catch (error) {
          console.error('Failed to load settings:', error);
          this.error = 'Failed to load settings';
        } finally {
          this.loading = false;
        }
      },
      settingLabel(key) {
        return String(key)
          .replace(/[._-]+/g, ' ')
          .replace(/\b\w/g, ch => ch.toUpperCase());
      },
      async saveSetting(setting) {
        this.savingKey = setting.key;
        try {
          const response = await fetch(this.urls.saveSetting, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              key: setting.key,
              value: setting.value
            })
          });
          const result = await response.json();

          if (result.status === 'success') {
            Notiflix.Notify.success('Setting saved: ' + setting.key, {
              position: 'right-top',
              timeout: 3000,
            });
          } else {
            Notiflix.Notify.failure(result.message || 'Failed to save setting', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        } catch (error) {
          console.error('Failed to save setting:', error);
          Notiflix.Notify.failure('Failed to save setting', {
            position: 'right-top',
            timeout: 3000,
          });
        } finally {
          this.savingKey = null;
        }
      }
    }
  }).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initSettingManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initSettingManager();
      }
    });
  }
});
