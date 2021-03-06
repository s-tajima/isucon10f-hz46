# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/api/config_change.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("google/api/config_change.proto", :syntax => :proto3) do
    add_message "google.api.ConfigChange" do
      optional :element, :string, 1
      optional :old_value, :string, 2
      optional :new_value, :string, 3
      optional :change_type, :enum, 4, "google.api.ChangeType"
      repeated :advices, :message, 5, "google.api.Advice"
    end
    add_message "google.api.Advice" do
      optional :description, :string, 2
    end
    add_enum "google.api.ChangeType" do
      value :CHANGE_TYPE_UNSPECIFIED, 0
      value :ADDED, 1
      value :REMOVED, 2
      value :MODIFIED, 3
    end
  end
end

module Google
  module Api
    ConfigChange = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("google.api.ConfigChange").msgclass
    Advice = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("google.api.Advice").msgclass
    ChangeType = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("google.api.ChangeType").enummodule
  end
end
