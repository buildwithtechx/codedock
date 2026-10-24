import { AppearanceSetting } from './appearance-setting';
import { LanguageSetting } from './language-setting';
import { PreferencesSetting } from './preferences-setting';

export function GeneralSettings() {
  return (
    <div className="space-y-6">
      <AppearanceSetting />
      <LanguageSetting />
      <PreferencesSetting />
    </div>
  );
}
