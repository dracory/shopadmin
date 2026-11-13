function initDiscountManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.discountManagerUrls || {},
        discounts: [],
        selectedDiscounts: [],
        selectAll: false,
        loading: false,
        error: null,
        showFilterModal: false,
        showCreateModal: false,
        creating: false,
        newDiscount: {
          title: ''
        },
        currentPage: 0,
        perPage: 10,
        totalDiscounts: 0,
        // Applied conditions (sent to the server)
        conditions: [],
        // Modal state — cloned from conditions when opening
        modalConditions: [],
        sortByColumn: 'created_at',
        sortOrder: 'desc',
        // Filter condition options — must match the Go handler fields.
        conditionOptions: [
          { value: 'code', label: 'Code', operators: ['equals'], inputType: 'text', placeholder: 'discount code' },
          { value: 'status', label: 'Status', operators: ['equals'], inputType: 'select', placeholder: '', options: [
            { value: 'active', label: 'Active' },
            { value: 'inactive', label: 'Inactive' },
            { value: 'draft', label: 'Draft' },
          ] },
          { value: 'type', label: 'Type', operators: ['equals'], inputType: 'select', placeholder: '', options: [
            { value: 'percent', label: 'Percent' },
            { value: 'amount', label: 'Amount' },
          ] },
          { value: 'id', label: 'ID', operators: ['equals'], inputType: 'text', placeholder: 'discount id' },
          { value: 'date_from', label: 'Created after', operators: ['equals'], inputType: 'date', placeholder: '' },
          { value: 'date_to', label: 'Created before', operators: ['equals'], inputType: 'date', placeholder: '' },
        ],
        opLabels: { equals: '=', contains: 'contains' }
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalDiscounts / this.perPage);
      },
      filterStatus() {
        if (this.conditions.length === 0) return 'Showing all discounts';
        return this.conditions.length + ' filter(s) active';
      },
      hasActiveFilters() {
        return this.conditions.length > 0;
      }
    },
    mounted() {
      const urlParams = new URLSearchParams(window.location.search);

      // Try JSON filters first (new format)
      const filtersRaw = urlParams.get('filters');
      if (filtersRaw) {
        try {
          const parsed = JSON.parse(filtersRaw);
          if (Array.isArray(parsed)) {
            this.conditions = parsed.filter(c => c && c.field && c.value);
          }
        } catch (e) { /* ignore malformed */ }
      } else {
        // Fallback: legacy individual params
        if (urlParams.get('status')) this.conditions.push({ field: 'status', operator: 'equals', value: urlParams.get('status') });
        if (urlParams.get('created_from')) this.conditions.push({ field: 'date_from', operator: 'equals', value: urlParams.get('created_from') });
        if (urlParams.get('created_to')) this.conditions.push({ field: 'date_to', operator: 'equals', value: urlParams.get('created_to') });
      }

      this.sortByColumn = urlParams.get('sort_by') || 'created_at';
      this.sortOrder = urlParams.get('sort_order') || 'desc';
      this.currentPage = parseInt(urlParams.get('page') || '0', 10);
      this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
      this.loadDiscounts();
    },
    methods: {
      async loadDiscounts() {
        this.loading = true;
        this.error = null;
        try {
          const response = await fetch(this.urls.loadDiscounts, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              page: this.currentPage,
              per_page: this.perPage,
              conditions: this.conditions,
              sort_by: this.sortByColumn,
              sort: this.sortOrder
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            if (result.data.urls) {
              this.urls = { ...this.urls, ...result.data.urls };
            }
            this.discounts = result.data.discounts;
            this.totalDiscounts = result.data.total || 0;
          } else {
            this.error = result.message || 'Failed to load discounts';
          }
        } catch (error) {
          console.error('Failed to load discounts:', error);
          this.error = 'Failed to load discounts';
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
        const maxPage = Math.ceil(this.totalDiscounts / this.perPage) - 1;
        if (page > maxPage) page = maxPage;
        this.currentPage = page;
        this.applyFilters();
      },
      changePerPage() {
        this.perPage = parseInt(this.perPage, 10);
        this.currentPage = 0;
        this.applyFilters();
      },
      openFilterModal() {
        this.modalConditions = this.conditions.map(c => ({ ...c }));
        if (this.modalConditions.length === 0) {
          this.modalConditions.push({ field: '', operator: 'equals', value: '' });
        }
        this.showFilterModal = true;
      },
      closeFilterModal() {
        this.showFilterModal = false;
      },
      addModalCondition() {
        this.modalConditions.push({ field: '', operator: 'equals', value: '' });
      },
      removeModalCondition(i) {
        this.modalConditions.splice(i, 1);
      },
      clearModalConditions() {
        this.modalConditions = [];
      },
      onFieldChange(c) {
        const opt = this.getFieldOpt(c.field);
        if (opt) {
          c.operator = opt.operators[0];
        }
        c.value = '';
      },
      applyConditions() {
        this.conditions = this.modalConditions.filter(c => c.field && String(c.value).trim());
        this.showFilterModal = false;
        this.currentPage = 0;
        this.applyFilters();
      },
      applyFilters() {
        const params = new URLSearchParams();
        const existing = new URLSearchParams(window.location.search);
        const controller = existing.get('controller');
        if (controller) params.set('controller', controller);

        if (this.conditions.length > 0) params.set('filters', JSON.stringify(this.conditions));
        params.set('page', this.currentPage);
        params.set('per_page', this.perPage);
        params.set('sort_order', this.sortOrder);
        params.set('sort_by', this.sortByColumn);

        const newUrl = `${window.location.pathname}?${params.toString()}`;
        window.history.pushState({}, '', newUrl);

        this.loadDiscounts();
      },
      clearFilters() {
        this.conditions = [];
        this.currentPage = 0;
        this.applyFilters();
      },
      // === Filter option helpers ===
      getFieldOpt(v) { return this.conditionOptions.find(o => o.value === v); },
      opLabel(op) { return this.opLabels[op] || op; },
      getFieldOperators(v) { const o = this.getFieldOpt(v); return o ? o.operators : ['equals']; },
      getFieldInputType(v) { const o = this.getFieldOpt(v); return o ? o.inputType : 'text'; },
      getFieldPlaceholder(v) { const o = this.getFieldOpt(v); return o ? o.placeholder : ''; },
      getFieldSelectOptions(v) { const o = this.getFieldOpt(v); return o && o.options ? o.options : []; },
      async createDiscount() {
        if (!this.newDiscount.title.trim()) {
          Notiflix.Notify.failure('Title is required', {
            position: 'right-top',
            timeout: 3000,
          });
          return;
        }

        this.creating = true;
        try {
          const response = await fetch(this.urls.createDiscount, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              title: this.newDiscount.title.trim()
            })
          });
          const result = await response.json();

          if (result.status === 'success') {
            Notiflix.Notify.success('Discount created successfully', {
              position: 'right-top',
              timeout: 3000,
            });
            this.showCreateModal = false;
            this.newDiscount.title = '';
            this.loadDiscounts();
          } else {
            Notiflix.Notify.failure(result.message || 'Failed to create discount', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        } catch (error) {
          console.error('Failed to create discount:', error);
          Notiflix.Notify.failure('Failed to create discount', {
            position: 'right-top',
            timeout: 3000,
          });
        } finally {
          this.creating = false;
        }
      },
      async deleteDiscount(discountId) {
        Swal.fire({
          title: 'Delete Discount',
          text: 'Are you sure you want to delete this discount?',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete it'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.deleteDiscount, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                discount_id: discountId
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.discounts = this.discounts.filter(d => d.id !== discountId);
              Swal.fire('Success', 'Discount deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete discount', 'error');
            }
          } catch (error) {
            console.error('Failed to delete discount:', error);
            Swal.fire('Error', 'Failed to delete discount', 'error');
          }
        });
      },
      async deleteSelectedDiscounts() {
        Swal.fire({
          title: 'Delete Discounts',
          text: `Are you sure you want to delete ${this.selectedDiscounts.length} discount(s)?`,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete them'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.deleteSelectedDiscounts, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                bulk_discount_ids: this.selectedDiscounts
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.discounts = this.discounts.filter(d => !this.selectedDiscounts.includes(d.id));
              this.selectedDiscounts = [];
              this.selectAll = false;
              Swal.fire('Success', 'Discounts deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete discounts', 'error');
            }
          } catch (error) {
            console.error('Failed to delete discounts:', error);
            Swal.fire('Error', 'Failed to delete discounts', 'error');
          }
        });
      },
      toggleSelectAll() {
        if (this.selectAll) {
          this.selectedDiscounts = this.discounts.map(d => d.id);
        } else {
          this.selectedDiscounts = [];
        }
      },
      getStatusBadgeClass(status) {
        switch (status) {
          case 'active':
            return 'bg-success';
          case 'inactive':
            return 'bg-secondary';
          case 'draft':
            return 'bg-warning';
          default:
            return 'bg-secondary';
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
    initDiscountManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initDiscountManager();
      }
    });
  }
});
