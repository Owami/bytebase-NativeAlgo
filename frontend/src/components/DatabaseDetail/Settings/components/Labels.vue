<template>
  <div class="space-y-4">
    <div class="flex items-center">
      <div class="flex-1">
        <div class="flex items-center">
          <p class="text-lg font-medium leading-7 text-main flex">
            {{ $t("common.labels") }}
          </p>
        </div>
      </div>
    </div>
    <div>
      <LabelListEditor
        ref="labelListEditorRef"
        v-model:kv-list="state.kvList"
        :readonly="state.mode === 'view'"
        :show-errors="dirty"
        class="max-w-[30rem]"
      />
    </div>
    <div class="flex flex-row justify-end items-center gap-x-3">
      <template v-if="state.mode === 'view'">
        <NButton :disabled="!allowAdmin" @click="beginEdit">
          {{ $t("common.edit") }}
        </NButton>
      </template>
      <template v-if="state.mode === 'edit'">
        <NButton @click="handleCancel">
          {{ $t("common.cancel") }}
        </NButton>
        <NButton
          :disabled="!allowSave || state.isUpdating"
          :loading="state.isUpdating"
          type="primary"
          @click="handleSave"
        >
          {{ $t("common.save") }}
        </NButton>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { cloneDeep, isEqual } from "lodash-es";
import { NButton } from "naive-ui";
import { computed, reactive, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { LabelListEditor } from "@/components/Label/";
import {
  pushNotification,
  useCurrentUserV1,
  useDatabaseV1Store,
} from "@/store";
import { type ComposedDatabase } from "@/types";
import { Database } from "@/types/proto/v1/database_service";
import {
  PRESET_LABEL_KEYS,
  convertKVListToLabels,
  convertLabelsToKVList,
  hasPermissionInProjectV1,
  hasWorkspacePermissionV1,
} from "@/utils";

type LocalState = {
  kvList: { key: string; value: string }[];
  mode: "view" | "edit";
  isUpdating: boolean;
};

const props = defineProps<{
  database: ComposedDatabase;
}>();

const { t } = useI18n();
const labelListEditorRef = ref<InstanceType<typeof LabelListEditor>>();
const me = useCurrentUserV1();
const state = reactive<LocalState>({
  kvList: [],
  mode: "view",
  isUpdating: false,
});

const allowAdmin = computed(() => {
  const project = props.database.projectEntity;
  return (
    hasWorkspacePermissionV1(
      "bb.permission.workspace.manage-label",
      me.value.userRole
    ) ||
    hasPermissionInProjectV1(
      project.iamPolicy,
      me.value,
      "bb.permission.project.manage-general"
    )
  );
});

const convert = () => {
  const labels = cloneDeep(props.database.labels);
  // Pre-fill preset label keys with empty values
  for (const key of PRESET_LABEL_KEYS) {
    if (!(key in labels)) {
      labels[key] = "";
    }
  }
  return convertLabelsToKVList(labels, true /* sort */);
};

const dirty = computed(() => {
  const original = convert();
  const local = state.kvList;
  return !isEqual(original, local);
});

const allowSave = computed(() => {
  if (!dirty.value) return false;
  const errors = labelListEditorRef.value?.flattenErrors ?? [];
  return errors.length === 0;
});

const beginEdit = () => {
  state.mode = "edit";
};

const handleCancel = () => {
  state.kvList = convert();
  state.mode = "view";
};

const handleSave = async () => {
  if (!allowSave.value) return;
  state.isUpdating = true;
  try {
    // Won't omit empty `tenant` value.
    // Otherwise the server API won't update the labels field correctly.
    const labels = convertKVListToLabels(state.kvList, false /* !omitEmpty */);

    const patch = {
      ...Database.fromPartial(props.database),
      labels,
    };
    await useDatabaseV1Store().updateDatabase({
      database: patch,
      updateMask: ["labels"],
    });
    pushNotification({
      module: "bytebase",
      style: "SUCCESS",
      title: t("common.updated"),
    });
    state.mode = "view";
  } finally {
    state.isUpdating = false;
  }
};

watch(
  allowAdmin,
  () => {
    state.mode = "view";
  },
  { immediate: true }
);

watch(
  () => props.database.labels,
  () => {
    state.kvList = convert();
  },
  {
    immediate: true,
    deep: true,
  }
);
</script>