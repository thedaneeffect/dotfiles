#!/bin/bash
# macOS Non-Default Settings
# These are `defaults write` commands that deviate from macOS defaults.
# Generated 2026-02-27

# ===========================================================================
# Global (NSGlobalDomain)
# ===========================================================================

# Dark mode
defaults write NSGlobalDomain AppleInterfaceStyle -string "Dark"

# 24-hour time
defaults write NSGlobalDomain AppleICUForce24HourTime -bool true

# Show all file extensions
defaults write NSGlobalDomain AppleShowAllExtensions -bool true

# Full keyboard access (tab through all controls)
defaults write NSGlobalDomain AppleKeyboardUIMode -int 1

# Fn key shows function keys (not emoji/special)
defaults write NSGlobalDomain com.apple.keyboard.fnState -bool true

# Disable sound feedback when adjusting volume
defaults write NSGlobalDomain com.apple.sound.beep.feedback -int 0

# Fast key repeat
defaults write NSGlobalDomain KeyRepeat -int 2

# Short delay until key repeat
defaults write NSGlobalDomain InitialKeyRepeat -int 15

# Disable auto-capitalization
defaults write NSGlobalDomain NSAutomaticCapitalizationEnabled -bool false

# Disable inline predictions
defaults write NSGlobalDomain NSAutomaticInlinePredictionEnabled -bool false

# Disable auto-period with double space
defaults write NSGlobalDomain NSAutomaticPeriodSubstitutionEnabled -bool false

# Disable auto-correct
defaults write NSGlobalDomain NSAutomaticSpellingCorrectionEnabled -bool false

# Disable web auto-correct
defaults write NSGlobalDomain WebAutomaticSpellingCorrectionEnabled -bool false

# Don't minimize on double-click title bar
defaults write NSGlobalDomain AppleMiniaturizeOnDoubleClick -bool false

# Trackpad scaling speed
defaults write NSGlobalDomain com.apple.trackpad.scaling -float 0.875

# Force click enabled
defaults write NSGlobalDomain com.apple.trackpad.forceClick -bool true

# Smart quotes style
defaults write NSGlobalDomain KB_DoubleQuoteOption -string "\u201cabc\u201d"
defaults write NSGlobalDomain KB_SingleQuoteOption -string "\u2018abc\u2019"

# ===========================================================================
# Dock
# ===========================================================================

# Auto-hide dock
defaults write com.apple.dock autohide -bool true

# Dock icon size
defaults write com.apple.dock tilesize -int 64

# Disable magnification
defaults write com.apple.dock magnification -bool false

# Enable app expose gesture
defaults write com.apple.dock showAppExposeGestureEnabled -bool true

# Hot corner: bottom-right = Quick Note (14)
defaults write com.apple.dock wvous-br-corner -int 14

# ===========================================================================
# Finder
# ===========================================================================

# Default view: column view
defaults write com.apple.finder FXPreferredViewStyle -string "clmv"

# Search current folder by default
defaults write com.apple.finder FXDefaultSearchScope -string "SCcf"

# Disable extension change warning
defaults write com.apple.finder FXEnableExtensionChangeWarning -bool false

# Auto-empty trash after 30 days
defaults write com.apple.finder FXRemoveOldTrashItems -bool true

# Don't warn before emptying trash
defaults write com.apple.finder WarnOnEmptyTrash -bool false

# Allow quitting Finder
defaults write com.apple.finder QuitMenuItem -bool true

# Show volumes on desktop
defaults write com.apple.finder ShowExternalHardDrivesOnDesktop -bool true
defaults write com.apple.finder ShowHardDrivesOnDesktop -bool true
defaults write com.apple.finder ShowMountedServersOnDesktop -bool true
defaults write com.apple.finder ShowRemovableMediaOnDesktop -bool true

# Arrange by date modified
defaults write com.apple.finder FK_ArrangeBy -string "Date Modified"

# New window target: home folder
defaults write com.apple.finder NewWindowTarget -string "PfCm"

# Disable iCloud Desktop & Documents
defaults write com.apple.finder FXICloudDriveDesktop -bool false
defaults write com.apple.finder FXICloudDriveDocuments -bool false

# ===========================================================================
# Screenshot
# ===========================================================================

# Disable shadow in screenshots
defaults write com.apple.screencapture disable-shadow -bool true

# Target: clipboard
defaults write com.apple.screencapture target -string "clipboard"

# Style: selection
defaults write com.apple.screencapture style -string "selection"

# Don't save screenshot selections
defaults write com.apple.screencapture save-selections -bool false

# Show clicks in recordings
defaults write com.apple.screencapture showsClicks -bool true

# ===========================================================================
# Spotlight
# ===========================================================================

# Disable web suggestions
defaults write com.apple.spotlight SuggestionsEnabled -int 0

# Disable Siri suggestions
defaults write com.apple.spotlight SiriSuggestionsEnabled -int 0

# ===========================================================================
# Privacy & Advertising
# ===========================================================================

# Disable personalized ads
defaults write com.apple.AdLib allowApplePersonalizedAdvertising -bool false

# Disable advertising identifier
defaults write com.apple.AdLib allowIdentifierForAdvertising -bool false

# Force limit ad tracking
defaults write com.apple.AdLib forceLimitAdTracking -bool true

# Disable lookup suggestions
defaults write com.apple.lookup.shared LookupSuggestionsDisabled -int 1

# ===========================================================================
# Safari
# ===========================================================================

# Disable universal search / search suggestions
defaults write com.apple.Safari UniversalSearchEnabled -bool false
defaults write com.apple.Safari SuppressSearchSuggestions -bool true

# Send Do Not Track header
defaults write com.apple.Safari SendDoNotTrackHTTPHeader -bool true

# Enable developer menu
defaults write com.apple.Safari IncludeDevelopMenu -bool true
defaults write com.apple.Safari WebKitDeveloperExtrasEnabledPreferenceKey -bool true
defaults write com.apple.Safari "WebKitPreferences.developerExtrasEnabled" -bool true

# Require authentication for private browsing
defaults write com.apple.Safari PrivateBrowsingRequiresAuthentication -bool true

# Enable extensions
defaults write com.apple.Safari ExtensionsEnabled -bool true

# ===========================================================================
# Window Manager
# ===========================================================================

# Hide desktop items (click desktop to reveal)
defaults write com.apple.WindowManager HideDesktop -bool true

# ===========================================================================
# Menu Bar Clock
# ===========================================================================

# Digital clock
defaults write com.apple.menuextra.clock IsAnalog -bool false

# Show AM/PM
defaults write com.apple.menuextra.clock ShowAMPM -bool true

# Show day of week
defaults write com.apple.menuextra.clock ShowDayOfWeek -bool true

# ===========================================================================
# Trackpad
# ===========================================================================

# Tap to click
defaults write com.apple.AppleMultitouchTrackpad Clicking -bool true
defaults write com.apple.driver.AppleBluetoothMultitouch.trackpad Clicking -bool true

# Scroll zoom with modifier (Ctrl)
defaults write com.apple.AppleMultitouchTrackpad HIDScrollZoomModifierMask -int 524288
defaults write com.apple.driver.AppleBluetoothMultitouch.trackpad HIDScrollZoomModifierMask -int 524288

# ===========================================================================
# Accessibility
# ===========================================================================

# Enable scroll gesture zoom (Ctrl + scroll)
defaults write com.apple.universalaccess closeViewScrollWheelToggle -bool true
defaults write com.apple.universalaccess closeViewScrollWheelModifiersInt -int 524288

# ===========================================================================
# Crash Reporter
# ===========================================================================

# Suppress crash reporter dialog
defaults write com.apple.CrashReporter DialogType -string "none"

# ===========================================================================
# Diagnostics & Analytics (System-level)
# ===========================================================================

# Disable sharing analytics with Apple
# (These are in /Library/Application Support/CrashReporter/DiagnosticMessagesHistory.plist)
# AutoSubmit = 0
# ThirdPartyDataSubmit = 0

# ===========================================================================
# Restart affected apps after applying
# ===========================================================================
killall Dock
killall Finder
killall SystemUIServer
killall cfprefsd
