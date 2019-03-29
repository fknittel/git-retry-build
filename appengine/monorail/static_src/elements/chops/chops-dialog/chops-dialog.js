// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import '@polymer/polymer/polymer-legacy.js';
import {PolymerElement, html} from '@polymer/polymer';
import {IronFocusablesHelper} from
  '@polymer/iron-overlay-behavior/iron-focusables-helper.js';
import {mixinBehaviors} from '@polymer/polymer/lib/legacy/class.js';

const ESC_KEYCODE = 27;
const TAB_KEYCODE = 9;

/**
 * `<chops-dialog>` displays a modal/dialog overlay.
 *
 * @customElement
 * @polymer
 */
export class ChopsDialog extends mixinBehaviors(
  [IronFocusablesHelper], PolymerElement) {
  static get template() {
    return html`
      <style>
        :host {
          position: fixed;
          z-index: 9999;
          left: 0;
          top: 0;
          width: 100%;
          height: 100%;
          overflow: auto;
          background-color: rgba(0,0,0,0.4);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        :host(:not([opened])), [hidden] {
          display: none;
          visibility: hidden;
        }
        :host([close-on-outside-click]),
        :host([close-on-outside-click]) .dialog::backdrop {
          /* TODO(zhangtiff): Deprecate custom backdrop in favor of native
           * browser backdrop.
           */
          cursor: pointer;
        }
        .dialog {
          background: none;
          border: 0;
          overflow: auto;
          max-width: 90%;
        }
        .dialog-content {
          /* This extra div is here because otherwise the browser can't
           * differentiate between a click event that hits the dialog element or
           * its backdrop pseudoelement.
           */
          box-sizing: border-box;
          background: white;
          padding: 1em 16px;
          cursor: default;
          box-shadow: 0px 3px 20px 0px hsla(0, 0%, 0%, 0.4);

          /*
           * Would have really preferred to just have the dialog itself be the
           * :host element so that users could style it more easily, but this
           * causes issues with the backdrop, because there's not a good way to
           * make a child element appear behind its parent without using wrappers.
           */
          @apply --chops-dialog-theme;
        }
      </style>
      <dialog class="dialog" role="dialog">
        <div class="dialog-content">
          <slot></slot>
        </div>
      </dialog>
    `;
  }

  static get is() {
    return 'chops-dialog';
  }

  static get properties() {
    return {
      /**
       * Whether the dialog should currently be displayed or not.
       */
      opened: {
        type: Boolean,
        value: false,
        reflectToAttribute: true,
        observer: '_openedChanged',
      },
      /**
       * A boolean that determines whether clicking outside of the dialog
       * window should close it.
       */
      closeOnOutsideClick: {
        type: Boolean,
        value: true,
        reflectToAttribute: true,
      },
      /**
       * A function fired when the element tries to change its own opened
       * state. This is useful if you want the dialog state managed outside
       * of the dialog instead of with internal state. (ie: with Redux)
       */
      onOpenedChange: Object,
      /**
       * Allow people to exit the dialog using keyboard shortcuts. Defaults
       * to the escape key.
       */
      exitKeys: {
        type: Array,
        value: [ESC_KEYCODE],
      },
      _boundKeydownHandler: Object,
      _previousFocusedElement: Object,
    };
  }

  ready() {
    super.ready();

    this._boundKeydownHandler = this._keydownHandler.bind(this);

    this.addEventListener('click', (evt) => {
      if (!this.opened || !this.closeOnOutsideClick) return;

      const hasDialog = evt.composedPath().find(
        (node) => {
          return node.classList && node.classList.contains('dialog-content');
        }
      );
      if (hasDialog) return;

      this.close();
    });
  }

  connectedCallback() {
    super.connectedCallback();
    window.addEventListener('keydown', this._boundKeydownHandler, true);
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    window.removeEventListener('keydown', this._boundKeydownHandler,
      true);
  }

  _keydownHandler(event) {
    if (!this.opened) return;
    let keyCode = event.keyCode;

    // Handle closing hot keys.
    if (this.exitKeys.includes(keyCode)) {
      this.close();
    }

    // For accessibility, prevent tabbing outside of the dialog.
    // Key code 9 is tab.
    if (keyCode === TAB_KEYCODE) {
      const tabbables = this.getTabbableNodes(this);
      const active = this._getActiveElement();

      if (!tabbables || !tabbables.length) {
        event.preventDefault();
      }

      if (event.shiftKey) {
        // backwards tab.
        if (active === tabbables[0]) {
          event.preventDefault();
          tabbables[tabbables.length - 1].focus();
        }
      } else {
        // forward tab.
        if (active === tabbables[tabbables.length - 1]) {
          event.preventDefault();
          tabbables[0].focus();
        }
      }
    }
  }

  close() {
    if (this.onOpenedChange) {
      this.onOpenedChange(false);
    } else {
      this.opened = false;
    }
  }

  open() {
    if (this.onOpenedChange) {
      this.onOpenedChange(true);
    } else {
      this.opened = true;
    }
  }

  toggle() {
    this.opened = !this.opened;
  }

  _getActiveElement() {
    // document.activeElement alone isn't sufficient to find the active
    // element within shadow dom.
    let active = document.activeElement || document.body;
    let activeRoot = active.shadowRoot || active.root;
    while (activeRoot && activeRoot.activeElement) {
      active = activeRoot.activeElement;
      activeRoot = active.shadowRoot || active.root;
    }
    return active;
  }

  _openedChanged(opened) {
    const dialog = this.shadowRoot.querySelector('dialog');
    if (opened) {
      // For accessibility, we want to ensure we remember the element that was
      // focused before this dialog opened.
      this._previousFocusedElement = this._getActiveElement();

      if (dialog.showModal) {
        dialog.showModal();
      } else {
        dialog.setAttribute('open', 'true');
      }

      // Focus the first element within the dialog when it's opened.
      const tabbables = this.getTabbableNodes(this);
      if (tabbables && tabbables.length) {
        tabbables[0].focus();
      } else if (this._previousFocusedElement) {
        this._previousFocusedElement.blur();
      }
    } else {
      if (dialog.close) {
        dialog.close();
      } else {
        dialog.setAttribute('open', undefined);
      }

      if (this._previousFocusedElement) {
        const element = this._previousFocusedElement;
        setTimeout(function() {
          // HACK. This is to prevent a possible accessibility bug where
          // using a keypress to trigger a button that exits a modal causes
          // the modal to immediately re-open because the button that
          // originally opened the modal refocuses, and the keypress
          // propagates.
          element.focus();
        }, 1);
      }
    }
  }
}
customElements.define(ChopsDialog.is, ChopsDialog);
