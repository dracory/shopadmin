function initSettingManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.settingManagerUrls || {},
        settings: [],
        selectedSettings: [],
        selectAll: false,
        loading: false,
        error: null,
        showEditModal: false,
        editingSetting: null,
        saving: false,
        form: { key: '', value: '' },
        currency: 'USD',
        savingCurrency: false,
        currencies: [
          'USD', 'EUR', 'GBP', 'JPY', 'AUD', 'CAD', 'CHF', 'CNY',
          'SEK', 'NZD', 'INR', 'BRL', 'ZAR', 'PLN', 'BGN', 'RON'
        ],
        currentPage: 0,
        perPage: 10,
        totalSettings: 0,
        filters: {
          key: ''
        },
        sortByColumn: 'setting_key',
        sortOrder: 'asc'
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalSettings / this.perPage);
      }
    },
    mounted() {
      const urlParams = new URLSearchParams(window.location.search);
      this.filters.key = urlParams.get('key') || '';
      this.sortByColumn = urlParams.get('sort_by') || 'setting_key';
      this.sortOrder = urlParams.get('sort_order') || 'asc';
      this.currentPage = parseInt(urlParams.get('page') || '0', 10);
      this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
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
              page: this.currentPage,
              per_page: this.perPage,
              key: this.filters.key,
              sort_by: this.sortByColumn,
              sort: this.sortOrder
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            if (result.data.urls) {
              this.urls = { ...this.urls, ...result.data.urls };
            }
            this.settings = result.data.settings;
            this.totalSettings = result.data.total || 0;
            this.currency = result.data.currency || 'USD';
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
      sortBy(column) {
        if (this.sortByColumn === column) {
          this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
        } else {
          this.sortByColumn = column;
          this.sortOrder = 'asc';
        }
        this.currentPage = 0;
        this.applyFilters();
      },
      goToPage(page) {
        if (page < 0) return;
        const maxPage = Math.ceil(this.totalSettings / this.perPage) - 1;
        if (page > maxPage) page = maxPage;
        this.currentPage = page;
        this.applyFilters();
      },
      changePerPage() {
        this.perPage = parseInt(this.perPage, 10);
        this.currentPage = 0;
        this.applyFilters();
      },
      applyFilters() {
        const params = new URLSearchParams();
        if (this.filters.key) params.set('key', this.filters.key);
        params.set('page', this.currentPage);
        params.set('per_page', this.perPage);
        params.set('sort_order', this.sortOrder);
        params.set('sort_by', this.sortByColumn);

        const newUrl = `${window.location.pathname}?controller=settings&${params.toString()}`;
        window.history.pushState({}, '', newUrl);

        this.loadSettings();
      },
      clearFilters() {
        this.filters.key = '';
        this.currentPage = 0;
        this.applyFilters();
      },
      openCreateModal() {
        this.editingSetting = null;
        this.form = { key: '', value: '' };
        this.showEditModal = true;
      },
      openEditModal(setting) {
        this.editingSetting = setting;
        this.form = { key: setting.key, value: setting.value };
        this.showEditModal = true;
      },
      async saveSetting() {
        if (!this.form.key.trim()) {
          Notiflix.Notify.failure('Key is required', {
            position: 'right-top',
            timeout: 3000,
          });
          return;
        }

        this.saving = true;
        try {
          const response = await fetch(this.urls.saveSetting, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              key: this.form.key.trim(),
              value: this.form.value
            })
          });
          const result = await response.json();

          if (result.status === 'success') {
            Notiflix.Notify.success('Setting saved successfully', {
              position: 'right-top',
              timeout: 3000,
            });
            this.showEditModal = false;
            this.loadSettings();
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
          this.saving = false;
        }
      },
      async saveCurrency() {
        this.savingCurrency = true;
        try {
          const response = await fetch(this.urls.saveSetting, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              key: 'currency',
              value: this.currency
            })
          });
          const result = await response.json();

          if (result.status === 'success') {
            Notiflix.Notify.success('Currency saved: ' + this.currency, {
              position: 'right-top',
              timeout: 3000,
            });
          } else {
            Notiflix.Notify.failure(result.message || 'Failed to save currency', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        } catch (error) {
          console.error('Failed to save currency:', error);
          Notiflix.Notify.failure('Failed to save currency', {
            position: 'right-top',
            timeout: 3000,
          });
        } finally {
          this.savingCurrency = false;
        }
      },
      async deleteSetting(settingId) {
        Swal.fire({
          title: 'Delete Setting',
          text: 'Are you sure you want to delete this setting?',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete it'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.deleteSetting, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                setting_id: settingId
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.settings = this.settings.filter(s => s.id !== settingId);
              Swal.fire('Success', 'Setting deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete setting', 'error');
            }
          } catch (error) {
            console.error('Failed to delete setting:', error);
            Swal.fire('Error', 'Failed to delete setting', 'error');
          }
        });
      },
      async deleteSelectedSettings() {
        Swal.fire({
          title: 'Delete Settings',
          text: `Are you sure you want to delete ${this.selectedSettings.length} setting(s)?`,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete them'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.deleteSelectedSettings, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                bulk_setting_ids: this.selectedSettings
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.settings = this.settings.filter(s => !this.selectedSettings.includes(s.id));
              this.selectedSettings = [];
              this.selectAll = false;
              Swal.fire('Success', 'Settings deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete settings', 'error');
            }
          } catch (error) {
            console.error('Failed to delete settings:', error);
            Swal.fire('Error', 'Failed to delete settings', 'error');
          }
        });
      },
      toggleSelectAll() {
        if (this.selectAll) {
          this.selectedSettings = this.settings.map(s => s.id);
        } else {
          this.selectedSettings = [];
        }
      },
      formatDate(dateString) {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString();
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
