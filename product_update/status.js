function initStatusApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initStatusApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                collapsed: false,
                status: 'draft',
                originalStatus: '',
                productID: statusProductId
            };
        },
        methods: {
            async loadStatus() {
                this.loading = true;
                try {
                    const response = await fetch(statusUrlLoad, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ product_id: this.productID })
                    });

                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        this.status = data.status_value || 'draft';
                        this.originalStatus = this.status;
                    } else {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    console.error('Error loading status:', error);
                    Notiflix.Notify.failure('Failed to load status: ' + error.message);
                } finally {
                    this.loading = false;
                }
            },

            async saveStatus() {
                if (!this.status) {
                    Notiflix.Notify.failure('Status is required');
                    return;
                }

                this.loading = true;
                try {
                    const response = await fetch(statusUrlSave, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            product_id: this.productID,
                            status: this.status
                        })
                    });

                    const data = await response.json();
                    if (data.status === 'success') {
                        Notiflix.Notify.success(data.message || 'Status saved successfully');
                        this.originalStatus = this.status;
                    } else {
                        Notiflix.Notify.failure(data.message || 'Failed to save status');
                    }
                } catch (error) {
                    console.error('Error saving status:', error);
                    Notiflix.Notify.failure('Failed to save status: ' + error.message);
                } finally {
                    this.loading = false;
                }
            }
        },
        mounted() {
            this.loadStatus();
        }
    }).mount('#status-app');
}

document.addEventListener('DOMContentLoaded', initStatusApp);
