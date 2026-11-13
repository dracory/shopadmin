function initAdminNotesApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initAdminNotesApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                collapsed: false,
                memo: '',
                originalMemo: '',
                productID: adminNotesProductId
            };
        },
        methods: {
            async loadNotes() {
                this.loading = true;
                try {
                    const response = await fetch(adminNotesUrlLoad, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ product_id: this.productID })
                    });

                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        this.memo = data.memo || '';
                        this.originalMemo = this.memo;
                    } else {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    console.error('Error loading admin notes:', error);
                    Notiflix.Notify.failure('Failed to load admin notes: ' + error.message);
                } finally {
                    this.loading = false;
                }
            },

            async saveNotes() {
                this.loading = true;
                try {
                    const response = await fetch(adminNotesUrlSave, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            product_id: this.productID,
                            memo: this.memo
                        })
                    });

                    const data = await response.json();
                    if (data.status === 'success') {
                        Notiflix.Notify.success(data.message || 'Admin notes saved successfully');
                        this.originalMemo = this.memo;
                    } else {
                        Notiflix.Notify.failure(data.message || 'Failed to save admin notes');
                    }
                } catch (error) {
                    console.error('Error saving admin notes:', error);
                    Notiflix.Notify.failure('Failed to save admin notes: ' + error.message);
                } finally {
                    this.loading = false;
                }
            }
        },
        mounted() {
            this.loadNotes();
        }
    }).mount('#admin-notes-app');
}

document.addEventListener('DOMContentLoaded', initAdminNotesApp);
