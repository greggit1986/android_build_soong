// Copyright 2020 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dexpreopt

import (
	"fmt"

	"android/soong/android"
)

type fakeToolBinary struct {
	android.ModuleBase
}

func (m *fakeToolBinary) GenerateAndroidBuildActions(ctx android.ModuleContext) {}

func (m *fakeToolBinary) HostToolPath() android.OptionalPath {
	return android.OptionalPathForPath(android.PathForTesting("dex2oat"))
}

func fakeToolBinaryFactory() android.Module {
	module := &fakeToolBinary{}
	android.InitAndroidArchModule(module, android.HostSupported, android.MultilibFirst)
	return module
}

func RegisterToolModulesForTest(ctx android.RegistrationContext) {
	ctx.RegisterModuleType("fake_tool_binary", fakeToolBinaryFactory)
}

func BpToolModulesForTest() string {
	return `
		fake_tool_binary {
			name: "dex2oatd",
		}
	`
}

func CompatLibDefinitionsForTest() string {
	bp := ""

	// For class loader context and <uses-library> tests.
	dexpreoptModules := []string{"android.test.runner"}
	dexpreoptModules = append(dexpreoptModules, CompatUsesLibs...)
	dexpreoptModules = append(dexpreoptModules, OptionalCompatUsesLibs...)

	for _, extra := range dexpreoptModules {
		bp += fmt.Sprintf(`
			java_library {
				name: "%s",
				srcs: ["a.java"],
				sdk_version: "none",
				system_modules: "stable-core-platform-api-stubs-system-modules",
				compile_dex: true,
				installable: true,
			}
		`, extra)
	}

	return bp
}

var PrepareForTestWithDexpreoptCompatLibs = android.GroupFixturePreparers(
	android.FixtureAddFile("defaults/dexpreopt/compat/a.java", nil),
	android.FixtureAddTextFile("defaults/dexpreopt/compat/Android.bp", CompatLibDefinitionsForTest()),
)

var PrepareForTestWithFakeDex2oatd = android.GroupFixturePreparers(
	android.FixtureRegisterWithContext(RegisterToolModulesForTest),
	android.FixtureAddTextFile("defaults/dexpreopt/Android.bp", BpToolModulesForTest()),
)

// Prepares a test fixture by enabling dexpreopt, registering the fake_tool_binary module type and
// using that to define the `dex2oatd` module.
var PrepareForTestByEnablingDexpreopt = android.GroupFixturePreparers(
	FixtureModifyGlobalConfig(func(android.PathContext, *GlobalConfig) {}),
)

// FixtureModifyGlobalConfig enables dexpreopt (unless modified by the mutator) and modifies the
// configuration.
func FixtureModifyGlobalConfig(configModifier func(ctx android.PathContext, dexpreoptConfig *GlobalConfig)) android.FixturePreparer {
	return android.FixtureModifyConfig(func(config android.Config) {
		// Initialize the dexpreopt GlobalConfig to an empty structure. This has no effect if it has
		// already been set.
		pathCtx := android.PathContextForTesting(config)
		dexpreoptConfig := GlobalConfigForTests(pathCtx)
		SetTestGlobalConfig(config, dexpreoptConfig)

		// Retrieve the existing configuration and modify it.
		dexpreoptConfig = GetGlobalConfig(pathCtx)
		configModifier(pathCtx, dexpreoptConfig)
	})
}

// FixtureSetArtBootJars enables dexpreopt and sets the ArtApexJars property.
func FixtureSetArtBootJars(bootJars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.ArtApexJars = android.CreateTestConfiguredJarList(bootJars)
	})
}

// FixtureSetBootJars enables dexpreopt and sets the BootJars property.
func FixtureSetBootJars(bootJars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.BootJars = android.CreateTestConfiguredJarList(bootJars)
	})
}

// FixtureSetApexBootJars sets the ApexBootJars property in the global config.
func FixtureSetApexBootJars(bootJars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.ApexBootJars = android.CreateTestConfiguredJarList(bootJars)
	})
}

// FixtureSetStandaloneSystemServerJars sets the StandaloneSystemServerJars property.
func FixtureSetStandaloneSystemServerJars(jars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.StandaloneSystemServerJars = android.CreateTestConfiguredJarList(jars)
	})
}

// FixtureSetSystemServerJars sets the SystemServerJars property.
func FixtureSetSystemServerJars(jars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.SystemServerJars = android.CreateTestConfiguredJarList(jars)
	})
}

// FixtureSetApexSystemServerJars sets the ApexSystemServerJars property in the global config.
func FixtureSetApexSystemServerJars(jars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.ApexSystemServerJars = android.CreateTestConfiguredJarList(jars)
	})
}

// FixtureSetApexStandaloneSystemServerJars sets the ApexStandaloneSystemServerJars property in the
// global config.
func FixtureSetApexStandaloneSystemServerJars(jars ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.ApexStandaloneSystemServerJars = android.CreateTestConfiguredJarList(jars)
	})
}

// FixtureSetPreoptWithUpdatableBcp sets the PreoptWithUpdatableBcp property in the global config.
func FixtureSetPreoptWithUpdatableBcp(value bool) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.PreoptWithUpdatableBcp = value
	})
}

// FixtureSetBootImageProfiles sets the BootImageProfiles property in the global config.
func FixtureSetBootImageProfiles(profiles ...string) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(ctx android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.BootImageProfiles = android.PathsForSource(ctx, profiles)
	})
}

// FixtureDisableGenerateProfile sets the DisableGenerateProfile property in the global config.
func FixtureDisableGenerateProfile(disable bool) android.FixturePreparer {
	return FixtureModifyGlobalConfig(func(_ android.PathContext, dexpreoptConfig *GlobalConfig) {
		dexpreoptConfig.DisableGenerateProfile = disable
	})
}
